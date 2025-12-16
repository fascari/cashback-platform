# Web3 Cashback Platform

Monorepo backend system that issues cashback as crypto tokens minted on a blockchain. Combines traditional backend components with Web3 concepts using event-driven architecture.

- **Off-chain**: Purchases, cashback calculation, business rules, persistence
- **On-chain**: Token minting, ownership, balance tracking

## Architecture

Follows **DDD** and **Clean Architecture**: `handler → usecase → repository`

- Usecases depend on interfaces declared in the usecase layer
- Repositories implement those interfaces
- Domain is pure (no framework dependencies)

## Repository Structure

```
cashback-platform/
├── .golangci.yml              # Linting configuration
├── go.work                    # Go workspace
├── Makefile                   # Build automation
│
├── services/
│   ├── cashback-service-api/  # REST entrypoint
│   ├── mint-consumer/         # Async event consumer
│   └── blockchain-adapter/    # gRPC service
│
├── proto/                     # Shared gRPC contracts
│   └── token.proto
│
├── db/
│   └── schema.sql            # Conceptual schemas
│
└── docs/
    ├── architecture.md
    ├── events.md
    └── specs/
```

## Services

### Cashback Service API
REST entry point for purchases and cashback rules. Persists state in PostgreSQL and publishes events (Outbox Pattern).

### Mint Consumer
Async event consumer with idempotent processing. Triggers token minting via gRPC.

### Blockchain Adapter
gRPC service that abstracts blockchain interaction.

## Technology Stack

- **Language**: Go 1.25+
- **Dependency Injection**: Uber Fx
- **Configuration**: Viper
- **HTTP Framework**: Chi
- **ORM**: GORM
- **Database**: PostgreSQL
- **Messaging**: NATS JetStream
- **RPC**: gRPC
- **Testing**: Mockery v2

## Getting Started

### Prerequisites

- Go 1.25+
- PostgreSQL
- NATS Server with JetStream enabled
- protoc (Protocol Buffers compiler)
- golangci-lint
- mockery v2

### Setup

```bash
make deps        # Download dependencies
make db-setup    # Setup databases
make proto       # Generate protobuf code
```

### Running

```bash
make run-blockchain-adapter  # Terminal 1
make run-mint-consumer       # Terminal 2
make run-cashback-service    # Terminal 3
```

### Building & Testing

```bash
make build  # Build all services
make test   # Run tests
make lint   # Lint code
```

## Development Guidelines

- **File Naming**: Lowercase, no underscores (except `*_test.go`)
- **Code Style**: Follow [Google Go Style Guide](https://google.github.io/styleguide/go/guide)
- **Architecture**: handler → usecase → repository, interfaces in usecase layer
- **Testing**: Use Mockery v2 for generating mocks from interfaces

## Make Commands

- `make all` - Download deps and build
- `make build/test/lint/fmt` - Build, test, lint, or format
- `make proto` - Generate protobuf code
- `make db-setup` - Create databases
- `make help` - Show all targets

## Documentation

- [Architecture](docs/architecture.md) - Detailed architecture documentation
- [Events](docs/events.md) - Event flow and descriptions
- [Refactoring Guide](REFACTORING.md) - Recent changes and guidelines
- [Project Spec](docs/specs/web3_cashback_project_spec.md) - Original specification

## Design Decisions

- Blockchain interaction isolated in dedicated adapter
- Async token minting with Outbox Pattern
- Blockchain as source of truth for balances
- Each service owns its PostgreSQL database

