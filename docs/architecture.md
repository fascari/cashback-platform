# Web3 Cashback Platform — Architecture

## 1. Overview

Backend system that issues cashback as **crypto tokens minted on a blockchain**. It integrates traditional backend components with Web3 concepts using an event-driven architecture and a clear separation between off-chain and on-chain responsibilities.

This repository is a **reference implementation** for exploring Web3, crypto, and blockchain concepts through a pragmatic, production-oriented design.

### Core Business Idea

Users earn cashback when making purchases. Instead of storing cashback only in a centralized database, the cashback value is represented as a **crypto token issued on a blockchain**, enabling transparent ownership, auditability, and clear separation of responsibilities.

### Off-chain vs On-chain

| Off-chain (PostgreSQL) | On-chain (Ethereum/Sepolia + Solana/Devnet) |
|---|---|
| Purchase handling | Token minting |
| Cashback calculation and business rules | Token ownership (wallet balances) |
| Event orchestration (Outbox Pattern) | Immutable audit trail (tx history) |
| Retry and failure handling | Deposit detection (Transfer events / signatures) |
| Off-chain ledger (cashback_ledger) | Confirmation tracking + chain reorg handling |

> Off-chain systems decide what should happen. On-chain systems guarantee the final and auditable state.

---

## 2. System Components

### cashback-service-api

REST entry point for purchases and cashback rules. Persists off-chain state in its own PostgreSQL database and publishes domain events via the Outbox Pattern. Has **no** direct dependency on blockchain libraries.

| Concern | Technology |
|---|---|
| HTTP | Chi |
| ORM | GORM |
| DI | Uber Fx |
| Config | Viper |

**Database** (`cashback_service_db`): `users`, `purchases`, `cashback_ledger`, `outbox_events`

### mint-consumer

Asynchronous event consumer. Processes `cashback.approved` events from NATS JetStream with idempotent handling and triggers token minting via gRPC.

| Concern | Technology |
|---|---|
| Messaging | NATS JetStream |
| gRPC Client | google.golang.org/grpc |
| DI | Uber Fx |
| Config | Viper |

**Database** (`mint_consumer_db`): `processed_events`, `mint_requests`

### blockchain-adapter

gRPC service that abstracts all blockchain interaction. Multi-chain adapter with a `ChainClient` interface enabling Ethereum and Solana support. Manages wallet signing, nonce/blockhash sequencing, deposit monitoring, transaction confirmation, and chain reorg detection.

| Concern | Technology |
|---|---|
| gRPC Server | google.golang.org/grpc |
| Ethereum | go-ethereum (ethclient, abigen bindings) |
| Solana | github.com/gagliardetto/solana-go |
| Distributed Lock | Redis (nonce serialization, fencing token) |
| DI | Uber Fx |
| Config | Viper |

**Database** (`blockchain_adapter_db`): `blockchain_transactions`, `wallet_nonces`, `detected_deposits`

---

## 3. Communication Patterns

| Boundary | Technology | Direction |
|---|---|---|
| External (client-facing) | REST | Client to cashback-service-api |
| Async internal | NATS JetStream | cashback-service-api to mint-consumer |
| Sync internal | gRPC | mint-consumer to blockchain-adapter |

---

## 4. Database Schema

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
| `processed_events` | Idempotency guard: stores `event_id` (UUID, unique) |
| `mint_requests` | Mint state machine: pending, processing, completed, failed |

### blockchain_adapter_db

| Table | Purpose |
|---|---|
| `blockchain_transactions` | On-chain tx lifecycle with `idempotency_key`, `transaction_hash`, `block_number`, `chain_id` |
| `wallet_nonces` | Per-wallet nonce counter; concurrency controlled by **Redis distributed lock + fencing token** (not SELECT FOR UPDATE — works correctly with multiple service replicas) |
| `detected_deposits` | Inbound deposits detected by the monitor: `chain_id`, `wallet_address`, `transaction_hash`, `amount`, `block_reference`, `status` |

Full schema: [`db/schema.sql`](../db/schema.sql)

---

## 5. Repository Structure

```text
cashback-platform/
├── contracts/                   # Solidity smart contracts
│   ├── CashbackToken.sol
│   ├── hardhat.config.ts
│   ├── package.json
│   └── scripts/deploy.ts
├── db/
│   └── schema.sql               # Conceptual schemas (all three DBs)
├── docs/
│   ├── architecture.md          # This file
│   ├── events.md                # Domain events, NATS config, Outbox Pattern
│   ├── decisions.md             # Technology stack, design decisions, trade-offs
│   └── diagrams/                # Mermaid diagrams, one .md per diagram
├── proto/
│   └── token.proto              # Shared gRPC contract
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
│           ├── chain/           # ChainClient interface + ChainRegistry
│           ├── config/
│           ├── contracts/       # Generated by abigen (Ethereum)
│           ├── domain/
│           ├── grpc/
│           ├── infra/
│           │   ├── ethereum/    # EthereumClient implements ChainClient
│           │   └── solana/      # SolanaClient implements ChainClient
│           ├── repository/
│           └── usecase/
├── .mise.toml                   # Task runner config
└── go.work                      # Go workspace
```

---

## 6. Diagrams

| # | Diagram |
|---|---|
| 1 | [High-Level System Flow](diagrams/01-high-level-system-flow.md) |
| 2 | [Off-chain vs On-chain Responsibilities](diagrams/02-offchain-vs-onchain.md) |
| 3 | [Domain Events Flow](diagrams/03-domain-events-flow.md) |
| 4 | [Blockchain Toolchain](diagrams/04-blockchain-toolchain.md) |
| 5 | [Transaction Lifecycle](diagrams/05-transaction-lifecycle.md) |
| 6 | [Wallet and Key Hierarchy](diagrams/06-wallet-key-hierarchy.md) |
| 7 | [Repository Structure](diagrams/07-repository-structure.md) |

---

## Further Reading

- [events.md](events.md): Domain events, NATS JetStream configuration, Outbox Pattern
- [decisions.md](decisions.md): Technology stack, design decisions, trade-offs, fault tolerance
