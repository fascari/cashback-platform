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

## 2. Compile the smart contract

```bash
cd contracts && npm install
mise run compile
```

Compiles `CashbackToken.sol` and writes artifacts to `contracts/artifacts/`.

## 3. Start a local Ethereum node

```bash
mise run evm
```

Starts a local EVM chain at `localhost:8545` and prints 20 pre-funded accounts with private keys. Keep this terminal open.

## 4. Deploy the contract

In a second terminal, from the `contracts/` directory:

```bash
mise run deploy
```

The output prints the deployed contract address. Copy it for the next step.

```text
Deploying with account: 0xf39Fd6e51...
CashbackToken deployed to: 0x5FbDB231...
```

## 5. Configure environment variables

Copy `.env.example` to `.env` in `services/blockchain-adapter/` and fill in:

```dotenv
ETHEREUM_RPC_URL=http://localhost:8545
CONTRACT_ADDRESS=0x<address from step 4>
WALLET_MNEMONIC=test test test test test test test test test test test junk
```

The mnemonic above is the Hardhat default. It matches the wallet that deployed the contract and holds owner permissions.

## 6. Generate Go bindings

From the repo root:

```bash
mise run contracts:bindings
```

Writes `services/blockchain-adapter/internal/contracts/cashbacktoken.go` from the compiled ABI. Re-run after any change to `CashbackToken.sol`.

## 7. Run the services

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
