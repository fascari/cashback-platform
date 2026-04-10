# mint-consumer

Asynchronous event consumer that processes `cashback.approved` events and triggers token minting via the blockchain adapter.

## Responsibilities

- Consume `cashback.approved` events from NATS JetStream
- Trigger token minting via gRPC to the blockchain-adapter service
- Track minting requests and their outcomes in a local database
- Retry failed mint operations via a scheduled job

## Events Consumed

| Subject | Description |
|---------|-------------|
| `cashback.approved` | Triggers token minting for the associated wallet |

## Architecture

```text
cmd/
├── main.go
└── modules/
    └── mint.go

internal/
├── app/
│   └── mint/
│       ├── domain/
│       ├── repository/
│       └── usecase/
│           ├── mintcashback/
│           └── retrymints/
├── bootstrap/
├── config/
├── consumer/
│   └── cashbackapproved/
└── infra/
    ├── database/
    ├── grpc/
    └── nats/
```

## Configuration

```env
APP_NAME=mint-consumer
APP_ENV=development

DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=mint_consumer_db
DATABASE_SSLMODE=disable

NATS_URL=nats://localhost:4222
BLOCKCHAIN_ADAPTER_GRPC_ADDRESS=localhost:50051
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

## Testing

```bash
go test ./...
```
