# cashback-service-api

HTTP API service for user, purchase, and cashback management. Publishes `cashback.approved` events via the Outbox Pattern for downstream token minting.

## Responsibilities

- Manage users and their wallet associations
- Record purchases and associate them with users
- Calculate and persist cashback for approved purchases
- Publish `cashback.approved` events reliably via the Outbox Pattern

## Architecture

The service follows clean architecture with one package per domain:

```text
cmd/api/
├── main.go
└── modules/
    ├── cashback.go
    ├── purchase.go
    ├── types.go
    └── user.go

internal/
├── app/
│   ├── cashback/
│   │   ├── domain/
│   │   ├── handler/{calculatecashback,findusercashback}/
│   │   ├── repository/
│   │   └── usecase/{calculatecashback,findusercashback}/
│   ├── purchase/
│   │   ├── domain/
│   │   ├── handler/{createpurchase,findpurchase}/
│   │   ├── repository/
│   │   └── usecase/{createpurchase,findpurchase}/
│   └── user/
│       ├── domain/
│       ├── handler/{createuser,finduser}/
│       ├── repository/
│       └── usecase/{createuser,finduser}/
├── bootstrap/
├── config/
├── database/
├── infra/
│   ├── grpc/
│   ├── nats/
│   └── messaging/outbox/
└── middleware/
```

Each domain exposes its behavior through use cases. Handlers parse HTTP requests and delegate to the use case. Repositories handle persistence and convert models to domain types.

The Outbox Pattern ensures `cashback.approved` is published atomically with the database write. A relay process polls the outbox table and forwards events to NATS JetStream.

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/users | Create a user |
| GET | /api/users/:id | Find a user by ID |
| POST | /api/purchases | Create a purchase |
| GET | /api/purchases/:id | Find a purchase by ID |
| POST | /api/cashback/calculate | Calculate and approve cashback |
| GET | /api/users/:user_id/cashback | Find cashback for a user |

## Configuration

```env
APP_NAME=cashback-service-api
APP_ENV=development
SERVER_PORT=8080

DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=cashback_service_db
DATABASE_SSLMODE=disable

NATS_URL=nats://localhost:4222
BLOCKCHAIN_ADAPTER_GRPC_ADDRESS=localhost:50051

OUTBOX_MAX_RETRIES=5
OUTBOX_POLL_INTERVAL_MS=100
```

## Running

```bash
mise run run
# or
go run cmd/api/main.go
```

Apply database migrations:

```bash
mise run db:migrate
```

## Testing

```bash
go test ./...
```
