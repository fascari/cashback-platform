# Local Development

This guide covers running the full cashback platform stack on a local machine.
No external accounts or testnets are required.

## Prerequisites

- Docker and Docker Compose
- Node.js 18+
- Go 1.25+
- [mise](https://mise.jdx.dev/) task runner

## 1. Start infrastructure

```bash
mise run up
```

Starts PostgreSQL (port 15432), NATS with JetStream, and Redis via Docker Compose.
Run `mise run db:setup` on first start to create the databases and apply schemas.

## 2. Compile the smart contract

```bash
cd contracts
npm install
mise run compile
```

Compiles `CashbackToken.sol` and writes artifacts to `contracts/artifacts/`.

## 3. Start a local Ethereum node

```bash
npx hardhat node
```

Starts a local EVM blockchain at `localhost:8545`.
The terminal prints 20 pre-funded accounts with their private keys. Keep this terminal open.

## 4. Deploy the contract

Open a second terminal in `contracts/`:

```bash
mise run deploy
```

Output:

```text
Deploying with account: 0xf39Fd6e51...
CashbackToken deployed to: 0x5FbDB231...
Set CONTRACT_ADDRESS=0x5FbDB231... in your .env
```

Copy the printed contract address.

## 5. Configure environment variables

Copy `.env.example` to `.env` in `services/blockchain-adapter/` and set:

```dotenv
ETHEREUM_RPC_URL=http://localhost:8545
CONTRACT_ADDRESS=0x<address from step 4>
WALLET_MNEMONIC=test test test test test test test test test test test junk
```

The mnemonic `test test...junk` is the Hardhat default. It corresponds to the
wallet that deployed the contract and holds owner permissions.

## 6. Generate Go bindings

```bash
# from repo root
mise run contracts:bindings
```

Writes `services/blockchain-adapter/internal/contracts/cashbacktoken.go` from
the compiled ABI. Re-run after any change to `CashbackToken.sol`.

## 7. Run the services

```bash
go run ./services/cashback-service-api/cmd/...
go run ./services/mint-consumer/cmd/...
go run ./services/blockchain-adapter/cmd/...
```

## Full cashback flow

```text
User purchase recorded
  -> cashback-service-api calculates amount
  -> publishes event to NATS
  -> mint-consumer consumes event
  -> calls blockchain-adapter via gRPC (MintToken)
  -> blockchain-adapter calls CashbackToken.mint() on the local node
  -> CBK tokens credited to the user wallet
```

## About Sepolia

Sepolia is a public Ethereum testnet. It is not required for local development.
Use Sepolia only when testing against a shared network before production.

To deploy to Sepolia:

1. Create a wallet in MetaMask.
2. Get test ETH from a Sepolia faucet (free).
3. Create an API key on [Infura](https://infura.io) or [Alchemy](https://alchemy.com).
4. Set the env vars and deploy:

```bash
export ETHEREUM_RPC_URL=https://sepolia.infura.io/v3/<your-key>
export DEPLOYER_PRIVATE_KEY=0x<your-wallet-private-key>
mise run deploy:sepolia
```
