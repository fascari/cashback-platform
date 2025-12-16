# Architecture Documentation

## Overview

The Web3 Cashback Platform is a monorepo backend system that issues cashback as crypto tokens minted on a blockchain. The architecture follows an event-driven design with clear separation between off-chain and on-chain responsibilities.

## System Components

### 1. Cashback Service API (`services/cashback-service-api`)

**Responsibility**: Core service handling purchases and cashback rules.

**Characteristics**:
- Acts as the system entry point
- Exposes REST API for external clients
- Persists off-chain state in its own PostgreSQL database
- Publishes domain events using the Outbox Pattern
- Does NOT depend on blockchain libraries directly

**Technology**:
- HTTP Framework: Chi
- ORM: GORM
- DI: Uber Fx
- Config: Viper

**Database Ownership**:
- `users` - User accounts and wallet addresses
- `purchases` - Purchase records
- `cashback_ledger` - Off-chain cashback representation
- `outbox_events` - Events pending publication

### 2. Mint Consumer (`services/mint-consumer`)

**Responsibility**: Asynchronous event processing and token minting orchestration.

**Characteristics**:
- Consumes events from NATS JetStream
- Ensures idempotent processing
- Triggers token minting via gRPC to Blockchain Adapter
- Publishes result events (success/failure)

**Technology**:
- Messaging: NATS JetStream
- gRPC Client: Connects to Blockchain Adapter
- DI: Uber Fx
- Config: Viper

**Database Ownership**:
- `processed_events` - Tracking for idempotency
- `mint_requests` - State of mint operations

### 3. Blockchain Adapter (`services/blockchain-adapter`)

**Responsibility**: Isolated abstraction for blockchain interaction.

**Characteristics**:
- Exposes gRPC interface
- Abstracts blockchain interaction
- Blockchain calls may be mocked or simulated
- Single point of contact with blockchain

**Technology**:
- gRPC Server
- DI: Uber Fx
- Config: Viper

**Database Ownership**:
- `blockchain_transactions` - Status of on-chain mint operations
- `wallet_nonces` - Nonce tracking for transactions

## Communication Patterns

### External Communication
- **REST API** exposed by Cashback Service API
- Used by external clients to create purchases

### Internal Communication
- **gRPC** between services
  - Mint Consumer → Blockchain Adapter (MintToken RPC)
  - Cashback Service API → Blockchain Adapter (GetBalance RPC, optional)
- **NATS JetStream** for event-driven communication
  - Cashback Service API → Mint Consumer (domain events)

## Data Flow

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │ REST
       ▼
┌─────────────────────────────┐
│   Cashback Service API      │
│  ┌───────────────────────┐  │
│  │ PostgreSQL (owned)    │  │
│  │ - users               │  │
│  │ - purchases           │  │
│  │ - cashback_ledger     │  │
│  │ - outbox_events       │  │
│  └───────────────────────┘  │
└──────┬──────────────────────┘
       │ NATS JetStream (events)
       ▼
┌─────────────────────────────┐
│      Mint Consumer          │
│  ┌───────────────────────┐  │
│  │ PostgreSQL (owned)    │  │
│  │ - processed_events    │  │
│  │ - mint_requests       │  │
│  └───────────────────────┘  │
└──────┬──────────────────────┘
       │ gRPC
       ▼
┌─────────────────────────────┐
│    Blockchain Adapter       │
│  ┌───────────────────────┐  │
│  │ PostgreSQL (owned)    │  │
│  │ - blockchain_txns     │  │
│  │ - wallet_nonces       │  │
│  └───────────────────────┘  │
└──────┬──────────────────────┘
       │ (mock/simulated)
       ▼
┌─────────────────────────────┐
│       Blockchain            │
└─────────────────────────────┘
```

## Data Ownership Principles

1. **No shared databases** - Each service owns its own PostgreSQL database
2. **No shared schema** - Each service manages its own schema and migrations
3. **Event-driven sync** - All cross-service data synchronization via events
4. **Single source of truth** - Blockchain is the source of truth for token balances

## Outbox Pattern

The Cashback Service implements the Outbox Pattern for reliable event publishing:

1. Business operation and event are written in the same database transaction
2. A separate process polls the outbox table
3. Events are published to NATS JetStream
4. Successfully published events are marked as processed

This ensures at-least-once delivery semantics.

## Fault Tolerance

### Idempotency
- All consumers are idempotent
- Events can be safely replayed
- Idempotency keys prevent duplicate processing

### Retry Mechanism
- Failed event processing is retried with exponential backoff
- Dead letter queues for events that exceed retry limits

### Fallback
- Temporary failures don't break the system flow
- Events remain in queues until successfully processed
- Circuit breaker pattern for blockchain calls

## Technology Decisions

| Concern | Technology | Rationale |
|---------|------------|-----------|
| Language | Go | Performance, concurrency, simplicity |
| DI | Uber Fx | Lifecycle management, testability |
| Config | Viper | Multi-source config, env vars support |
| HTTP | Chi | Lightweight, stdlib compatible |
| ORM | GORM | Productivity, migrations support |
| Messaging | NATS JetStream | Persistence, replay, exactly-once |
| RPC | gRPC | Type safety, performance, streaming |
| Database | PostgreSQL | ACID, reliability, JSONB support |

