# Mint Consumer

Asynchronous event consumer that processes cashback events and triggers token minting.

## Responsibilities

- Consume events from NATS JetStream
- Ensure idempotent processing
- Trigger token minting via gRPC to Blockchain Adapter
- Publish result events (success/failure)

## Events Consumed

- `cashback.approved` - Triggers token minting

## Events Produced

- `token.mint.requested` - When minting is initiated
- `token.minted` - When minting succeeds
- `token.mint.failed` - When minting fails

## Configuration

Environment variables:

```
APP_NAME=mint-consumer
APP_ENV=development
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=mint_consumer_db
NATS_URL=nats://localhost:4222
BLOCKCHAIN_ADAPTER_GRPC_ADDRESS=localhost:50051
```

## Running

```bash
go run cmd/main.go
```

## Project Structure

```
mint-consumer/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── domain/
│   │   ├── mint_request.go
│   │   └── processed_event.go
│   ├── repository/
│   │   ├── mint_request_repository.go
│   │   └── processed_event_repository.go
│   ├── service/
│   │   └── mint_service.go
│   ├── consumer/
│   │   └── cashback_consumer.go
│   └── infrastructure/
│       ├── database/
│       │   └── postgres.go
│       ├── nats/
│       │   └── client.go
│       └── grpc/
│           └── blockchain_client.go
├── go.mod
└── README.md
```

