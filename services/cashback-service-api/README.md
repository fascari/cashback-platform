# Cashback Service API

The core service handling purchases and cashback rules. Acts as the system entry point for the Web3 Cashback Platform.

## 🎯 Responsibilities

- Expose REST API for external clients
- Handle purchase registration
- **Calculate and approve cashback**
- Persist off-chain state in PostgreSQL
- **Publish domain events using the Outbox Pattern**
- Coordinate with Blockchain Adapter for token minting

---

## 🏗️ Architecture

This service follows **Domain-Driven Design (DDD)** principles with clear separation of concerns:

```
internal/
├── app/                    # Business domains
│   ├── user/              ✅ User management
│   ├── purchase/          ✅ Purchase tracking
│   └── cashback/          ✅ Cashback calculation & tracking
├── infrastructure/        # External adapters
│   ├── grpc/             # Blockchain Adapter client
│   ├── nats/             # Event streaming
│   └── outbox/           # Outbox pattern implementation
└── database/             # Database connection
```

### Key Design Patterns

- **Repository Pattern**: Data access abstraction
- **Use Case Pattern**: Application logic encapsulation
- **Outbox Pattern**: Reliable event publishing
- **Dependency Injection**: Via Uber Fx

---

## 📡 API Endpoints

### Users

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/users` | Register a new user |
| GET | `/api/users/:id` | Get user by ID |

### Purchases

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/purchases` | Create a new purchase |
| GET | `/api/purchases/:id` | Get purchase by ID |

### Cashback ⭐ NEW

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/cashback/calculate` | Calculate cashback for a purchase |
| GET | `/api/users/:user_id/cashback` | Get cashback summary for a user |

---

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- NATS JetStream
- Blockchain Adapter (gRPC service)

### Installation

```bash
# Clone repository
cd services/cashback-service-api

# Install dependencies
go mod download

# Run database migrations
make migrate-up

# Start service
go run cmd/api/main.go
```

### Using Docker Compose

```bash
# Start all dependencies
docker-compose up -d postgres nats

# Run service
make run
```

---

## ⚙️ Configuration

Environment variables:

```bash
# Application
APP_NAME=cashback-service-api
APP_ENV=development
SERVER_PORT=8080

# Database
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=cashback_service_db
DATABASE_SSLMODE=disable

# NATS
NATS_URL=nats://localhost:4222

# Blockchain Adapter
BLOCKCHAIN_ADAPTER_GRPC_ADDRESS=localhost:50051
```

---

## 📊 Database Schema

### Tables

- **users**: User accounts with wallet addresses
- **purchases**: Purchase records
- **cashback_ledger**: Off-chain cashback tracking
- **outbox_events**: Events pending publication

### Migrations

```bash
make migrate-up      # Apply migrations
make migrate-down    # Rollback migrations
make migrate-create  # Create new migration
```

---

## 🔄 Event Flow

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │ POST /cashback/calculate
       ▼
┌─────────────┐
│  API Layer  │
└──────┬──────┘
       │
       ▼
┌─────────────────────┐
│  Calculate UseCase  │
│  - Validate         │
│  - Calculate 5%     │
│  - Approve          │
│  - Persist          │
└──────┬──────────────┘
       │
       ▼
┌─────────────────────┐
│  Outbox Publisher   │
│  cashback.approved  │
└──────┬──────────────┘
       │
       ▼
┌─────────────┐
│ NATS Stream │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│Mint Consumer│ (separate service)
└─────────────┘
```

---

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific package
go test ./internal/app/cashback/...
```

### Example API Test

```bash
# 1. Create user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "external_id": "user123",
    "email": "user@example.com",
    "wallet_address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
  }'

# 2. Create purchase
curl -X POST http://localhost:8080/api/purchases \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "<USER_ID>",
    "amount": 100.00,
    "merchant": "Amazon"
  }'

# 3. Calculate cashback (5% of 100 = 5.00)
curl -X POST http://localhost:8080/api/cashback/calculate \
  -H "Content-Type: application/json" \
  -d '{
    "purchase_id": "<PURCHASE_ID>"
  }'

# 4. Get user cashback
curl http://localhost:8080/api/users/<USER_ID>/cashback
```

---

## 📚 Documentation

- [Cashback Module Documentation](./CASHBACK_MODULE_DOCUMENTATION.md) - Complete module reference
- [API Examples](./API_EXAMPLES.md) - Request/response examples
- [Refactoring Report](./REFACTORING_REPORT.md) - Architecture decisions
- [Project Spec](../../docs/specs/web3_cashback_platform_spec.md) - Business requirements

---

## 🏛️ Domain Events

### Published Events

| Event | Trigger | Consumer |
|-------|---------|----------|
| `purchase.created` | New purchase registered | N/A (future) |
| `cashback.approved` | Cashback calculated and approved | Mint Consumer |

### Event Schema: cashback.approved

```json
{
  "cashback_id": "uuid",
  "user_id": "uuid",
  "wallet_address": "0x...",
  "purchase_id": "uuid",
  "amount": 5.0,
  "cashback_percent": 5.0
}
```

---

## 🛠️ Development

### Project Structure

```
cmd/api/
├── main.go           # Application entry point
└── modules/          # Fx modules
    ├── user.go
    ├── purchase.go
    └── cashback.go

internal/
├── app/
│   └── cashback/
│       ├── domain/       # Business entities
│       ├── repository/   # Data access
│       ├── usecase/      # Business logic
│       └── handler/      # HTTP handlers
├── infrastructure/
│   └── outbox/          # Event publishing
└── bootstrap/           # App initialization
```

### Adding a New Feature

1. Define domain entity in `domain/`
2. Implement repository in `repository/`
3. Create use case in `usecase/`
4. Add HTTP handler in `handler/`
5. Register in Fx module

---

## 🐛 Troubleshooting

### Database Connection Failed

```bash
# Check PostgreSQL is running
docker ps | grep postgres

# Test connection
psql -h localhost -U postgres -d cashback_service_db
```

### NATS Connection Failed

```bash
# Check NATS is running
docker ps | grep nats

# Test connection
nats-cli server info
```

### Events Not Publishing

```bash
# Check outbox table
SELECT * FROM outbox_events WHERE published = false;

# Check publisher logs
make logs
```

---

## 📝 License

MIT License - see LICENSE file for details

---

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

---

## 📞 Support

For questions or issues, please open an issue on GitHub.

---

**Status**: ✅ Production Ready  
**Version**: 1.0.0  
**Last Updated**: December 16, 2025
```

## Project Structure

```
cashback-service-api/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── domain/
│   │   ├── user.go
│   │   ├── purchase.go
│   │   └── cashback.go
│   ├── repository/
│   │   ├── user_repository.go
│   │   ├── purchase_repository.go
│   │   ├── cashback_repository.go
│   │   └── outbox_repository.go
│   ├── service/
│   │   ├── user_service.go
│   │   ├── purchase_service.go
│   │   └── cashback_service.go
│   ├── handler/
│   │   ├── user_handler.go
│   │   ├── purchase_handler.go
│   │   └── cashback_handler.go
│   ├── outbox/
│   │   └── publisher.go
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

