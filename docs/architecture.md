# Web3 Cashback Platform — Architecture

## 1. Overview

Backend system that issues cashback as **crypto tokens minted on a blockchain**. It integrates traditional backend components with Web3 concepts using an event-driven architecture and a clear separation between off-chain and on-chain responsibilities.

This repository serves as a **reference implementation** and a **validation of Web3, crypto, and blockchain concepts** through a pragmatic, production-oriented design.

### Core Business Idea

Users earn cashback when making purchases. Instead of storing cashback only in a centralized database, the cashback value is represented as a **crypto token issued on a blockchain**, enabling transparent ownership, auditability, and clear separation of responsibilities.

### Off-chain vs On-chain

| Off-chain (PostgreSQL) | On-chain (Ethereum / Sepolia) |
|---|---|
| Purchase handling | Token minting |
| Cashback calculation and business rules | Token ownership (wallet balances) |
| Event orchestration (Outbox Pattern) | Immutable audit trail (tx history) |
| Retry and failure handling | |
| Off-chain ledger (cashback_ledger) | |

> **Off-chain systems decide what should happen. On-chain systems guarantee the final and auditable state.**

---

## 2. System Components

### cashback-service-api

REST entry point for purchases and cashback rules. Persists off-chain state in its own PostgreSQL database and publishes domain events via the Outbox Pattern. Does NOT depend on blockchain libraries directly.

| Concern | Technology |
|---|---|
| HTTP | Chi |
| ORM | GORM |
| DI | Uber Fx |
| Config | Viper |

**Database ownership** (`cashback_service_db`):
- `users` — User accounts and wallet addresses
- `purchases` — Purchase records
- `cashback_ledger` — Off-chain cashback representation
- `outbox_events` — Events pending publication

### mint-consumer

Asynchronous event consumer that processes `cashback.approved` events from NATS JetStream, ensures idempotent processing, and triggers token minting via gRPC to the blockchain-adapter.

| Concern | Technology |
|---|---|
| Messaging | NATS JetStream |
| gRPC Client | google.golang.org/grpc |
| DI | Uber Fx |
| Config | Viper |

**Database ownership** (`mint_consumer_db`):
- `processed_events` — Tracking for idempotency
- `mint_requests` — Mint state machine: pending, processing, completed, failed

### blockchain-adapter

gRPC service that abstracts all blockchain interaction. Single point of contact with the Ethereum network. Manages wallet signing, nonce sequencing, and transaction lifecycle.

| Concern | Technology |
|---|---|
| gRPC Server | google.golang.org/grpc |
| Ethereum | go-ethereum (ethclient, abigen bindings) |
| DI | Uber Fx |
| Config | Viper |

**Database ownership** (`blockchain_adapter_db`):
- `blockchain_transactions` — Status of on-chain mint operations
- `wallet_nonces` — Nonce tracking for concurrent transaction signing

---

## 3. Communication Patterns

### External

- **REST API** exposed by cashback-service-api (client-facing)

### Internal

- **NATS JetStream** for event-driven async communication (cashback-service-api to mint-consumer)
- **gRPC** for synchronous service-to-service calls (mint-consumer to blockchain-adapter)

---

## 4. Domain Events

### Event Catalog

| Event | Producer | Consumer | Trigger |
|---|---|---|---|
| `purchase.created` | cashback-service-api | cashback-service-api (internal) | Client submits a purchase via REST |
| `cashback.approved` | cashback-service-api (outbox) | mint-consumer | Cashback calculated and approved |
| `token.mint.requested` | mint-consumer | mint-consumer (internal state) | mint-consumer receives `cashback.approved` |
| `token.minted` | mint-consumer | cashback-service-api (optional) | blockchain-adapter confirms successful mint |
| `token.mint.failed` | mint-consumer | mint-consumer (retry logic) | blockchain-adapter returns an error |

### Event Payloads

#### cashback.approved

```json
{
  "event_id": "uuid",
  "event_type": "cashback.approved",
  "timestamp": "2024-01-15T10:30:01Z",
  "data": {
    "cashback_id": "uuid",
    "purchase_id": "uuid",
    "user_id": "uuid",
    "wallet_address": "0x...",
    "cashback_amount": 1.50,
    "token_amount": "1500000000000000000",
    "calculation_basis": {
      "purchase_amount": 150.00,
      "cashback_percentage": 1.0
    }
  }
}
```

#### token.minted

```json
{
  "event_id": "uuid",
  "event_type": "token.minted",
  "timestamp": "2024-01-15T10:30:05Z",
  "data": {
    "mint_request_id": "uuid",
    "cashback_id": "uuid",
    "user_id": "uuid",
    "wallet_address": "0x...",
    "token_amount": "1500000000000000000",
    "transaction_hash": "0x...",
    "block_number": 12345678,
    "minted_at": "2024-01-15T10:30:05Z"
  }
}
```

#### token.mint.failed

```json
{
  "event_id": "uuid",
  "event_type": "token.mint.failed",
  "timestamp": "2024-01-15T10:30:05Z",
  "data": {
    "mint_request_id": "uuid",
    "cashback_id": "uuid",
    "wallet_address": "0x...",
    "token_amount": "1500000000000000000",
    "error_code": "BLOCKCHAIN_UNAVAILABLE",
    "error_message": "Failed to connect to blockchain node",
    "retry_count": 1,
    "max_retries": 5,
    "next_retry_at": "2024-01-15T10:31:05Z"
  }
}
```

---

## 5. NATS JetStream Configuration

### Streams

| Stream | Subjects | Retention | MaxAge | Storage |
|---|---|---|---|---|
| `CASHBACK_EVENTS` | `cashback.>` | Limits | 7 days | File |
| `TOKEN_EVENTS` | `token.>` | Limits | 7 days | File |
| `PURCHASE_EVENTS` | `purchase.>` | Limits | 7 days | File |

### Consumers

| Consumer | Stream | Filter Subject | AckPolicy | MaxDeliver | AckWait |
|---|---|---|---|---|---|
| `mint-consumer` | CASHBACK_EVENTS | `cashback.approved` | Explicit | 5 | 30s |
| `cashback-service-token-updates` | TOKEN_EVENTS | `token.minted` | Explicit | 5 | 30s |

---

## 6. Outbox Pattern

The cashback-service-api uses the Outbox Pattern for reliable event publishing:

1. Business operation and event are written in the **same database transaction**
2. A separate goroutine polls the `outbox_events` table every 100ms
3. Events are published to NATS JetStream
4. Successfully published events are marked as `published`

This guarantees **at-least-once delivery** semantics — events are never lost even if NATS is temporarily unavailable.

---

## 7. Database Schema

Each service owns its own PostgreSQL database. There is **no shared database**.

### cashback_service_db

| Table | Purpose |
|---|---|
| `users` | User accounts with `wallet_address` (VARCHAR 42) |
| `purchases` | Purchase records with `amount`, `merchant_id`, `status` |
| `cashback_ledger` | Off-chain cashback with `token_amount` (wei as VARCHAR 78) |
| `outbox_events` | Pending events with `status` (pending / published / failed) |

### mint_consumer_db

| Table | Purpose |
|---|---|
| `processed_events` | Idempotency guard — stores `event_id` (UUID, unique) |
| `mint_requests` | Mint state machine: pending, processing, completed, failed |

### blockchain_adapter_db

| Table | Purpose |
|---|---|
| `blockchain_transactions` | On-chain tx lifecycle with `idempotency_key`, `transaction_hash`, `block_number` |
| `wallet_nonces` | Per-wallet nonce counter with `SELECT FOR UPDATE` for concurrency |

Full schema: [`db/schema.sql`](../db/schema.sql)

---

## 8. Fault Tolerance

### Idempotency

- All consumers are idempotent via `event_id` tracking in `processed_events`
- `idempotency_key` in `blockchain_transactions` prevents duplicate on-chain mints
- Events can be safely replayed

### Retry Strategy

Failed operations follow exponential backoff:

| Retry | Delay |
|---|---|
| 1 | 1s |
| 2 | 2s |
| 3 | 4s |
| 4 | 8s |
| 5 | 16s |

After 5 retries, events are moved to a dead letter queue for manual inspection.

### Fallback

- Temporary failures do not break the system flow
- Events remain in queues until successfully processed
- Circuit breaker pattern for blockchain RPC calls

---

## 9. Technology Stack

| Concern | Technology | Rationale |
|---|---|---|
| Language | Go 1.26 | Performance, concurrency, simplicity |
| DI | Uber Fx | Lifecycle management, testability |
| Config | Viper | Multi-source config, env vars support |
| HTTP | Chi | Lightweight, stdlib compatible |
| ORM | GORM | Productivity, migrations support |
| Messaging | NATS JetStream | Persistence, replay, at-least-once |
| RPC | gRPC + protobuf | Type safety, performance, streaming |
| Database | PostgreSQL | ACID, reliability, JSONB support |
| Blockchain | go-ethereum | Official Go Ethereum client |
| Smart Contract | Solidity + OpenZeppelin | ERC-20 standard, battle-tested |
| Contract Tooling | Hardhat | Compilation, local network, deploy scripts |
| Task Runner | Mise | Polyglot task runner, replaces Makefile |
| Testing | Mockery v2 | Mock generation from interfaces |

---

## 10. Design Decisions

| Decision | Rationale |
|---|---|
| Blockchain interaction isolated in dedicated adapter | Single point of change; other services stay blockchain-agnostic |
| Async consumers for blockchain operations | On-chain calls are slow and unreliable; async decouples availability |
| Outbox Pattern for event publishing | Guarantees at-least-once delivery without distributed transactions |
| Blockchain as source of truth for balances | Off-chain ledger is a cache; on-chain state is the authority |
| REST for external API | Simple, widely supported, no browser limitations |
| gRPC for internal communication | Type-safe, fast, streaming support |
| One database per service | No shared state; services communicate only via events or gRPC |

---

## 11. Trade-offs

- Blockchain is not used for business rules — avoids cost and latency
- Event-driven flow increases complexity but improves resilience and decoupling
- Off-chain ledger is a read cache, not the source of truth
- gRPC simplifies internal communication but limits direct browser access

---

## 12. Repository Structure

```
cashback-platform/
├── contracts/                   # Solidity smart contracts
│   ├── CashbackToken.sol
│   ├── hardhat.config.ts
│   ├── package.json
│   └── scripts/deploy.ts
├── db/
│   └── schema.sql               # Conceptual schemas (all 3 DBs)
├── docs/
│   ├── architecture.md           # This file
│   └── diagrams/                 # PlantUML diagrams (.puml)
├── proto/
│   └── token.proto               # Shared gRPC contract
├── services/
│   ├── cashback-service-api/
│   │   ├── cmd/api/main.go
│   │   └── internal/
│   │       ├── app/{domain}/
│   │       │   ├── domain/
│   │       │   ├── handler/{operation}/
│   │       │   ├── repository/
│   │       │   └── usecase/{operation}/
│   │       ├── bootstrap/
│   │       ├── config/
│   │       └── infra/messaging/outbox/
│   ├── mint-consumer/
│   │   ├── cmd/main.go
│   │   └── internal/
│   │       ├── consumer/
│   │       ├── domain/
│   │       ├── repository/
│   │       └── usecase/
│   └── blockchain-adapter/
│       ├── cmd/main.go
│       └── internal/
│           ├── config/
│           ├── contracts/         # Generated by abigen
│           ├── domain/
│           ├── grpc/
│           ├── infra/ethereum/
│           ├── repository/
│           └── usecase/
├── mise.toml                     # Task runner config
└── go.work                       # Go workspace
```

---

## Diagrams

Visual architecture diagrams in PlantUML format ([`docs/diagrams/`](./diagrams/)). Open `.puml` files in IntelliJ with the PlantUML plugin for live rendering.

| File | Description |
|---|---|
| `01-high-level-system-flow.puml` | How a purchase becomes an on-chain token |
| `02-offchain-vs-onchain.puml` | Off-chain vs On-chain responsibilities |
| `03-domain-events-flow.puml` | Sequence diagram of the full event flow |
| `04-blockchain-toolchain.puml` | Solidity to Go bindings pipeline |
| `05-transaction-lifecycle.puml` | State machine: pending to confirmed/failed |
| `06-wallet-key-hierarchy.puml` | Mnemonic to address derivation |
| `07-repository-structure.puml` | Where each concept lives in the codebase |

