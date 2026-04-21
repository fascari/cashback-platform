# Local Development

Runs the full stack on a local machine. No external accounts or testnets required.

## Prerequisites

- [mise](https://mise.jdx.dev/) task runner
- Go 1.26.1+
- Docker and Docker Compose
- Node.js 18+ (for Hardhat)

---

## 1. Start infrastructure

```bash
mise run up
```

Starts PostgreSQL on port 15432, NATS with JetStream, Redis, Jaeger, and the OTel Collector. Confirm all containers are healthy:

```bash
docker ps --format "table {{.Names}}\t{{.Status}}"
```

Expected: `postgres`, `nats`, and `redis` with status `(healthy)`.

## 2. Apply database migrations

On first run, create the databases and apply all schemas:

```bash
mise run db:setup
```

Expected: `no change` or migration applied messages per schema (`cashback`, `mint`, `blockchain`).

## 3. Provision NATS streams

```bash
mise run nats:setup
```

Creates four JetStream streams: `PURCHASE_EVENTS`, `CASHBACK_EVENTS`, `TOKEN_EVENTS`, and `DEPOSIT_EVENTS`. The command is idempotent. Running `run:api` or `run:consumer` via mise runs this step automatically.

## 4. Compile the smart contract

```bash
cd contracts && npm install
mise run compile
```

Compiles `CashbackToken.sol` and writes artifacts to `contracts/artifacts/`.

## 5. Start a local Ethereum node

```bash
mise run evm
```

Starts a Hardhat Network node at `localhost:8545` and prints 20 pre-funded accounts with private keys. Keep this terminal open. Expected output:

```text
Started HTTP and WebSocket JSON-RPC server at http://127.0.0.1:8545/
```

**Alternative:** [Anvil](https://book.getfoundry.sh/anvil/) from the Foundry suite is a lighter option. It exposes the same JSON-RPC interface and uses the same deterministic test accounts. Install Foundry and run:

```bash
anvil --chain-id 31337
```

The e2e test suite uses Anvil automatically via Docker. No manual setup is required for e2e runs.

## 6. Deploy the contract

In a second terminal, from the `contracts/` directory:

```bash
mise run deploy
```

The output prints the deployed contract address. Copy it for the next step.

```text
Deploying with account: 0xf39Fd6e51...
CashbackToken deployed to: 0x5FbDB231...
```

## 7. Configure environment variables

Copy `.env.example` to `.env` in `services/blockchain-adapter/` and fill in:

```dotenv
ETHEREUM_RPC_URL=http://localhost:8545
CONTRACT_ADDRESS=0x<address from step 6>
WALLET_MNEMONIC=test test test test test test test test test test test junk
```

The mnemonic is the Hardhat default. It matches the wallet that deployed the contract and holds owner permissions.

## 8. Generate Go bindings

From the repo root:

```bash
mise run contracts:bindings
```

Writes `services/blockchain-adapter/internal/contracts/cashbacktoken.go` from the compiled ABI. Re-run after any change to `CashbackToken.sol`.

## 9. Run the services

Each service runs in its own terminal:

```bash
mise run run:adapter   # blockchain-adapter    :50051
mise run run:api       # cashback-service-api  :8080
mise run run:consumer  # mint-consumer
```

Start the adapter first. The API creates the NATS streams on startup, which the consumer depends on.

Expected logs:

```text
gRPC server starting on port 50051   # blockchain-adapter
server started on port 8080          # cashback-service-api
cashback consumer started            # mint-consumer
```

---

## Running the cashback flow

Import the Postman collection:

```text
services/cashback-service-api/docs/api/postman_collection.json
```

The collection uses a `base_url` variable set to `http://localhost:8080` and auto-captures `user_id`, `purchase_id`, and `cashback_id` from responses via test scripts.

Run requests in this order:

1. **Users > Create User** — creates a user with a wallet address. The `user_id` is saved automatically to the collection variable.
2. **Purchases > Create Purchase** — registers a purchase for the user. The `purchase_id` is saved automatically.
3. **Cashback > Calculate Cashback** — calculates 5% cashback on the purchase amount. The record is saved with `status: approved` and a `cashback.approved` event is published to NATS via the transactional outbox.

After step 3, the consumer receives the event and logs:

```text
minting cashback cashback_id=X wallet=0x... amount=...
mint completed mint_id=X tx_hash=0x... block=X
```

On a duplicate message (idempotency check):

```text
mint skipped: duplicate event event_id=...
```

Verify the result:

- **Cashback > Find User Cashback** — confirms the cashback record with `status: approved`.
- **Users > Get User Balance** — returns the on-chain token balance in wei and the human-readable amount (e.g., `5.00` tokens).

To check the balance directly via gRPC:

```bash
grpcurl -plaintext \
  -d '{"wallet_address":"0xYourWallet"}' \
  localhost:50051 token.TokenService/GetBalance
```

---

## What happens end-to-end

**User and purchase creation** go through the REST API and are persisted in PostgreSQL under the `cashback` schema. No side effects occur at this stage.

**Cashback calculation** is the first critical step. When `POST /cashbacks/calculate/{purchase_id}` is called, the use case calculates 5% of the purchase amount and creates a cashback record with `status: approved`. Inside the same database transaction, it writes a `cashback.approved` event to the `cashback.outbox_events` table. This guarantees the event is never lost, even if the process crashes before publishing to NATS.

**Outbox relay** runs as a background goroutine inside cashback-service-api. It polls the outbox table every 5 seconds, picks up pending events, publishes them to the NATS JetStream stream `CASHBACK_EVENTS` on subject `cashback.approved`, and marks each event as `published`. The relay uses a transactional outbox pattern to ensure at-least-once delivery without dual-write risk.

**mint-consumer** receives the message from NATS using a durable pull consumer named `mint-consumer`. It deserializes the payload and calls the `mintcashback` use case, which first checks idempotency: if the `event_id` (a UUID) was already processed, the message is acknowledged and skipped. Otherwise, it creates a `mint_request` record and calls the blockchain-adapter via gRPC.

**blockchain-adapter** receives the `MintToken` gRPC call with the wallet address and token amount in wei. It resolves the nonce for the signing wallet, builds the transaction, signs it with the configured wallet (derived from the mnemonic in `.env`), and submits it to the EVM node. On the local Hardhat network, the transaction confirms immediately. The adapter returns the transaction hash and block number.

**mint-consumer** receives the result, marks the mint as `completed` in the `mint.mint_requests` table, and logs:

```text
mint completed mint_id=X tx_hash=0x... block=X
```

**The token balance** on the recipient wallet increases by the minted amount. This can be verified via `GET /users/{id}/balance`, which calls the blockchain-adapter's `Balance` gRPC method and returns the balance in wei alongside a human-readable token amount.

**Idempotency** is enforced at two levels: the outbox uses the cashback ID as the aggregate key, and the consumer stores each processed `event_id` in a `processed_events` table before creating the mint request. Replaying the same NATS message produces no duplicate mint.

**Retry** is triggered when the gRPC call fails. The mint is marked `failed` with an error code and the consumer retries it every 30 seconds up to 5 attempts. After 5 failures the record stays in `failed` state for inspection.

---

## Deposit detection flow

1. A `Transfer` event is emitted on the `CashbackToken` contract (from any source: mint, transfer, or external wallet)
2. `blockchain-adapter` polls the Ethereum node every 12 seconds for new Transfer events in the latest block range
3. Each detected event is saved to `blockchain.detected_deposits` with its transaction hash, wallet address, token amount, and block number
4. A `deposit.detected` event is published to the `DEPOSIT_EVENTS` NATS stream for downstream consumers

The monitor resumes from the last stored block number on restart, so no events are lost or reprocessed after a service restart.

### Testing deposit detection locally

With the Hardhat node running and the contract deployed, trigger a `Transfer` event by calling `mint()` from the deployer account using `cast` from Foundry:

```bash
cast send \
  --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
  --rpc-url http://127.0.0.1:8545 \
  <CONTRACT_ADDRESS> \
  "mint(address,uint256)" \
  0x70997970C51812dc3A010C7d01b50e0d17dc79C8 \
  1000000000000000000
```

The blockchain-adapter polls every 12 seconds. Within one cycle, a row appears in `blockchain.detected_deposits` and a `deposit.detected` event is published to the `DEPOSIT_EVENTS` NATS stream.

Verify:

```bash
psql "$POSTGRES_DSN_BLOCKCHAIN" \
  -c "SELECT transaction_hash, wallet_address, token_amount, block_number, status FROM detected_deposits ORDER BY id DESC LIMIT 5;"
```

---

## Retry behavior

- Retry loop runs every 30 seconds
- NATS consumer is configured with `MaxDeliver: 5`
- Failed mints are tracked in the `mint.mint_requests` table with `retry_count` and `error_message`

---

## Observability

Jaeger and the OpenTelemetry Collector start with the default stack — no extra step needed.

Tracing is enabled by default via `OTEL_ENABLED=true` in `.env.example`. To disable it for a run:

```bash
OTEL_ENABLED=false mise run run:api
```

### Viewing traces

Open [http://localhost:16686](http://localhost:16686). After making any HTTP request to `cashback-service-api`, select the service from the dropdown and click **Find Traces**.

### Pipeline

```
cashback-service-api
      | OTLP/gRPC :4317
      v
otel-collector  -->  Jaeger :4317 (internal)
                           |
                     UI :16686
```

The collector receives spans from the service, batches them, and forwards to Jaeger over OTLP.

### Verifying

```bash
docker compose ps jaeger otel-collector
curl http://localhost:16686/api/services
```

After the API starts with tracing enabled, `cashback-service-api` appears in the service list.

---

## E2E tests

The e2e suite runs the three services against a dedicated isolated stack (separate Postgres on port 25432, NATS on port 4322, Redis on port 6479) and exercises the full cashback and deposit detection flows through real HTTP calls and on-chain transactions.

### Running the full suite

```bash
mise run test:e2e
```

The script `scripts/run-e2e.sh` handles the full lifecycle:

1. Starts isolated Docker Compose infrastructure (`docker-compose.e2e.yml`)
2. Waits for Postgres, NATS, and Anvil
3. Deploys `CashbackToken` to the local Anvil node
4. Creates NATS streams via `cmd/nats-setup`
5. Applies DB migrations against the e2e databases
6. Builds and starts the three services
7. Waits for the API to be healthy on port 18080
8. Runs the test suite with `-tags=e2e`
9. Tears down all infrastructure

### Anvil EVM node

The e2e infrastructure includes an [Anvil](https://book.getfoundry.sh/anvil/) node (`ghcr.io/foundry-rs/foundry`) running on port `8545`. Anvil is the EVM simulator from the Foundry suite. It provides the same interface as Hardhat Network but ships as a standalone binary with no Node.js dependency.

The node starts with `--chain-id 31337`, matching `ETHEREUM_CHAIN_ID` in `.env`. The deployer account is the standard Foundry test account derived from the `"test test...junk"` mnemonic, which is also the `WALLET_MNEMONIC` in `.env`. After the node is ready, the script deploys `CashbackToken` via Hardhat, producing the deterministic address `0x5FbDB2315678afecb367f032d93F642f64180aa3` already set as `CONTRACT_ADDRESS` in `.env`.

### Blockchain tests

The e2e suite includes two tests that require a running EVM node.

`TestCashbackFlow_ShouldIncrementBalanceAfterMint` exercises the mint path: a cashback approval triggers a `CashbackToken.mint()` call and the test verifies the user balance returned by the API increases.

`TestDepositFlow_ShouldDetectOnChainTransfer` exercises the deposit detection path: the test mints tokens directly on the contract and verifies that `blockchain-adapter` stores the detected `Transfer` event in `detected_deposits` within 45 seconds.

Both tests run automatically as part of `mise run test:e2e` because Anvil is part of the e2e compose. If no EVM node is reachable on `ETHEREUM_RPC_URL`, both tests skip gracefully.

---

## Teardown

Stop service terminals with `Ctrl+C`, then:

```bash
mise run down
```

---

## Deploying to Sepolia

Sepolia is a public Ethereum testnet, not required for local development. Use it only when testing against a shared network.

1. Create a wallet in MetaMask and fund it from a Sepolia faucet.
2. Create an API key on [Infura](https://infura.io) or [Alchemy](https://alchemy.com).
3. Set the environment variables and deploy:

```bash
export ETHEREUM_RPC_URL=https://sepolia.infura.io/v3/<your-key>
export DEPLOYER_PRIVATE_KEY=0x<your-wallet-private-key>
mise run deploy:sepolia
```

