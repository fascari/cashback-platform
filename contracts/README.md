# Cashback Contracts

Solidity smart contracts for the cashback platform. Built with Hardhat.

## Contracts

### CashbackToken (CBK)

ERC-20 token with restricted mint and burn operations.

| Function | Access | Description |
|---|---|---|
| `mint(address to, uint256 amount)` | Owner only | Creates tokens and sends to `to` |
| `burn(address from, uint256 amount)` | Owner only | Destroys tokens from `from` |
| Standard ERC-20 | Public | `transfer`, `approve`, `allowance`, etc. |

The deployer wallet becomes the owner. Only the owner can mint or burn. Any other caller gets the transaction reverted.

`amount` is always in wei (1 token = `1e18` wei).

## Prerequisites

- Node.js 18+
- npm

## Setup

```bash
cd contracts
npm install
```

## Compile

```bash
# inside contracts/
mise run compile

# from the repo root
mise run contracts:compile
```

Generates artifacts in `artifacts/` and TypeScript bindings in `typechain-types/`.

## Generate Go bindings

Run after compiling. Requires `abigen` (installed automatically by the task).

```bash
# inside contracts/
mise run bindings

# from the repo root
mise run contracts:bindings
```

Writes `services/blockchain-adapter/internal/contracts/cashbacktoken.go`.
Re-run whenever `CashbackToken.sol` changes.

## Deploy

### Local (Hardhat node)

Start a local EVM node in one terminal:

```bash
npx hardhat node
```

Deploy in another:

```bash
npm run deploy:local
```

### Sepolia testnet

Set the required env vars:

```bash
export ETHEREUM_RPC_URL=https://sepolia.infura.io/v3/<your-key>
export DEPLOYER_PRIVATE_KEY=0x<private-key-of-deployer-wallet>
```

Deploy:

```bash
npm run deploy:sepolia
```

The deploy script prints the contract address. Set it in the service `.env`:

```dotenv
CONTRACT_ADDRESS=0x<printed-address>
```

## How it connects to blockchain-adapter

The Go service calls `mint()` via the generated binding when a cashback event is processed:

```text
mint-consumer → gRPC MintToken → blockchain-adapter → CashbackToken.mint()
```

The binding at `services/blockchain-adapter/internal/contracts/cashbacktoken.go`
is generated from the compiled ABI and is not committed to version control.
Regenerate with `mise run contracts:bindings` after any contract change.
