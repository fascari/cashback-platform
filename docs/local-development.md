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

Starts PostgreSQL on port 15432, NATS with JetStream, and Redis. On first run, create the databases and apply schemas:

```bash
mise run db:setup
```

## 2. Provision NATS streams

```bash
mise run nats:setup
```

Creates the three JetStream streams (`PURCHASE_EVENTS`, `CASHBACK_EVENTS`, `TOKEN_EVENTS`). The command is idempotent and exits immediately. `run:api` and `run:consumer` depend on this task, so running them directly via mise handles it automatically.

## 3. Compile the smart contract

```bash
cd contracts && npm install
mise run compile
```

Compiles `CashbackToken.sol` and writes artifacts to `contracts/artifacts/`.

## 4. Start a local Ethereum node

```bash
mise run evm
```

Starts a local EVM chain at `localhost:8545` and prints 20 pre-funded accounts with private keys. Keep this terminal open.

## 5. Deploy the contract

In a second terminal, from the `contracts/` directory:

```bash
mise run deploy
```

The output prints the deployed contract address. Copy it for the next step.

```text
Deploying with account: 0xf39Fd6e51...
CashbackToken deployed to: 0x5FbDB231...
```

## 6. Configure environment variables

Copy `.env.example` to `.env` in `services/blockchain-adapter/` and fill in:

```dotenv
ETHEREUM_RPC_URL=http://localhost:8545
CONTRACT_ADDRESS=0x<address from step 4>
WALLET_MNEMONIC=test test test test test test test test test test test junk
```

The mnemonic above is the Hardhat default. It matches the wallet that deployed the contract and holds owner permissions.

## 7. Generate Go bindings

From the repo root:

```bash
mise run contracts:bindings
```

Writes `services/blockchain-adapter/internal/contracts/cashbacktoken.go` from the compiled ABI. Re-run after any change to `CashbackToken.sol`.

## 8. Run the services

Each service runs in its own terminal:

```bash
mise run run:api       # cashback-service-api  :8080
mise run run:consumer  # mint-consumer
mise run run:adapter   # blockchain-adapter    :50051
```

---

## Full cashback flow

1. Purchase recorded via cashback-service-api
2. Cashback calculated and event published to NATS
3. mint-consumer picks up the event and calls blockchain-adapter via gRPC
4. blockchain-adapter signs and submits `CashbackToken.mint()` on the local node
5. CBK tokens credited to the user wallet

---

## E2E Tests

The e2e suite runs the three services against a dedicated isolated stack (separate Postgres on port 25432, NATS on port 4322, Redis on port 6479) and exercises the full cashback flow through real HTTP calls.

### Running the full suite

```bash
mise run test:e2e
```

The script `scripts/run-e2e.sh` handles the full lifecycle:

1. Starts isolated Docker Compose infrastructure (`docker-compose.e2e.yml`)
2. Waits for Postgres and NATS
3. Creates NATS streams via `cmd/nats-setup`
4. Applies DB migrations against the e2e databases
5. Starts the three services natively with `go run`
6. Waits for the API to be healthy on port 18080
7. Runs the test suite with `-tags=e2e`
8. Tears down all infrastructure

### Blockchain test

The test `TestCashbackFlow_ShouldIncrementBalanceAfterMint` requires a running Hardhat node. It is skipped by default. To include it:

```bash
# In one terminal
mise run evm

# Then run the e2e suite with blockchain enabled
E2E_BLOCKCHAIN=true mise run test:e2e
```

---

## Deploying to Sepolia

Sepolia is a public Ethereum testnet and is not required for local development. Use it only when testing against a shared network.

1. Create a wallet in MetaMask and fund it from a Sepolia faucet.
2. Create an API key on [Infura](https://infura.io) or [Alchemy](https://alchemy.com).
3. Set the environment variables and deploy:

```bash
export ETHEREUM_RPC_URL=https://sepolia.infura.io/v3/<your-key>
export DEPLOYER_PRIVATE_KEY=0x<your-wallet-private-key>
mise run deploy:sepolia
```

