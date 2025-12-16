# Blockchain Adapter

Isolated service that abstracts blockchain interaction and exposes a gRPC interface.

## Responsibilities

- Expose gRPC interface for token operations
- Abstract blockchain interaction
- Handle transaction submission and tracking
- Provide idempotent mint operations

## gRPC Services

- `MintToken` - Mint tokens to a wallet address
- `GetBalance` - Get token balance for a wallet
- `GetTransaction` - Get transaction status

## Configuration

Environment variables:

```
APP_NAME=blockchain-adapter
APP_ENV=development
GRPC_PORT=50051
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=blockchain_adapter_db
```

## Running

```bash
go run cmd/main.go
```

## Project Structure

```
blockchain-adapter/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── domain/
│   │   ├── blockchain_transaction.go
│   │   └── wallet_nonce.go
│   ├── repository/
│   │   ├── transaction_repository.go
│   │   └── nonce_repository.go
│   ├── service/
│   │   └── token_service.go
│   ├── grpc/
│   │   └── server.go
│   └── infrastructure/
│       └── database/
│           └── postgres.go
├── go.mod
└── README.md
```

## Notes

- Blockchain calls are mocked/simulated in this implementation
- The adapter provides idempotency via idempotency keys
- Transaction status is tracked in the local database

