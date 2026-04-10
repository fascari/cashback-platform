# blockchain-adapter

gRPC service that abstracts Ethereum blockchain interaction. Exposes token operations to internal services and handles nonce management, transaction submission, and idempotent retry logic.

## Responsibilities

- Expose a gRPC interface for token minting operations
- Submit and track transactions on the Ethereum network
- Manage wallet nonces with Redis-backed distributed locking
- Provide idempotent operations via idempotency keys

## gRPC Services

| RPC | Description |
|-----|-------------|
| `MintToken` | Mint tokens to a wallet address |
| `GetBalance` | Get token balance for a wallet |
| `GetTransaction` | Get transaction status by ID |

## Architecture

```text
cmd/
├── main.go
└── modules/
    └── token.go

internal/
├── bootstrap/
├── config/
├── contracts/
│   └── cashbacktoken.go
├── domain/
├── grpc/
│   └── server.go
├── infra/
│   ├── database/
│   └── redis/
├── repository/
│   ├── nonce/
│   └── transaction/
└── usecase/
```

## Configuration

```env
APP_NAME=blockchain-adapter
APP_ENV=development
GRPC_PORT=50051

POSTGRES_DSN_BLOCKCHAIN=postgres://cashback_app:cashback_app@localhost:15432/blockchain_adapter_db?sslmode=disable&search_path=blockchain

ETHEREUM_RPC_URL=https://sepolia.infura.io/v3/YOUR_PROJECT_ID
ETHEREUM_CHAIN_ID=11155111
CONTRACT_ADDRESS=0x...

WALLET_MNEMONIC=word1 word2 ... word12
WALLET_DERIVATION_PATH=m/44'/60'/0'/0/0

REDIS_URL=redis://localhost:6379
```

## Running

```bash
mise run run
# or
go run cmd/main.go
```

Apply database migrations:

```bash
mise run db:migrate
```

Regenerate Go contract bindings after ABI changes:

```bash
mise run contracts:bindings
```

## Testing

```bash
go test ./...
```
