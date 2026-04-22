# Local Setup

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

Alternative: [Anvil](https://book.getfoundry.sh/anvil/) from the Foundry suite is a lighter option. It exposes the same JSON-RPC interface and uses the same deterministic test accounts. Install Foundry and run:

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


## Teardown

Stop service terminals with `Ctrl+C`, then:

```bash
mise run down
```
