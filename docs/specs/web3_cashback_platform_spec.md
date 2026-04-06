# Web3 Cashback Platform - Project Specification

## 1. Project Overview

This document specifies a backend system that issues cashback as crypto tokens minted on a blockchain.

The system integrates traditional backend components with Web3 concepts, using an event-driven architecture and a clear separation between off-chain and on-chain responsibilities.

This repository validates and explores business, Web3, crypto, and blockchain concepts through a pragmatic and production-oriented design.

---

## 2. Core Business Idea

Users earn cashback when making purchases.

Instead of storing cashback only in a centralized database, the cashback value is represented as a crypto token issued on a blockchain, enabling:

- transparent ownership
- auditability
- clear separation of responsibilities

---

## 3. Web3, Blockchain, and Crypto Concepts

### Blockchain
A distributed and immutable ledger used as the source of truth for token ownership.

### Crypto Token
A digital asset issued on a blockchain, representing the cashback value earned by users. Implemented as a Solidity smart contract (`CashbackToken`) deployed on Ethereum.

### Web3
An application model where users directly own digital assets, while backend systems act as coordinators rather than custodians.

---

## 4. On-chain vs Off-chain Responsibilities

### Off-chain
- Purchase handling
- Cashback calculation
- Business rules
- Persistence
- Event orchestration
- Retries and failure handling

### On-chain
- Token minting
- Token ownership
- Final balance tracking

Off-chain systems decide what should happen. On-chain systems guarantee the final and auditable state.

---

## 5. System Architecture Components

- **Cashback Service**
  Core service responsible for purchases, cashback rules, and off-chain persistence. Exposes a REST API.

- **Mint Consumer**
  Asynchronous component that subscribes to `cashback.approved` events and triggers token minting via the Blockchain Adapter.

- **Blockchain Adapter**
  Isolated service that abstracts Ethereum interactions via gRPC. Manages nonce coordination and transaction lifecycle.

- **Smart Contract**
  Solidity contract (`CashbackToken`) managed with Hardhat. Defines the token minting logic executed on-chain.

- **Event Broker**
  NATS JetStream used for reliable and durable event delivery.

- **Database**
  PostgreSQL used for off-chain state (one database per service).

- **Cache**
  Redis used by the Blockchain Adapter for distributed nonce management.

- **Shared Kit**
  Internal Go module (`kit/`) with shared packages: `apperror`, `gormtx`, `logger`, `nats`, `ethereum`.

---

## 6. Database Schema (Per Service)

Each service manages its own database schema.

### cashback-service-api (`cashback` schema)

- **users**: users and associated wallet addresses
- **purchases**: purchase records
- **cashback_ledger**: off-chain representation of generated cashback (statuses: `pending`, `approved`, `minting`, `minted`, `failed`)
- **outbox_events**: pending events to be published via the Outbox Pattern

### mint-consumer (`mint` schema)

- **mint_requests**: mint operations with status tracking and retry metadata (statuses: `pending`, `processing`, `completed`, `failed`)
- **processed_events**: deduplication log for idempotent event consumption

### blockchain-adapter (`blockchain` schema)

- **blockchain_transactions**: on-chain transaction records with hash, block number, gas, and status (statuses: `pending`, `submitted`, `confirmed`, `failed`)
- **wallet_nonces**: per-wallet nonce tracking with optimistic locking via fence tokens
- **detected_deposits**: on-chain deposit events detected from the blockchain

---

## 7. Domain Events

The system publishes one domain event over NATS JetStream:

- **cashback.approved**: cashback amount has been calculated and approved; consumed by the Mint Consumer to trigger minting

Downstream operations (gRPC call to Blockchain Adapter, transaction submission, confirmation) are tracked as internal state transitions, not as published events.

### Fault Tolerance

- Retry with configurable backoff
- Replay via NATS JetStream durable consumers
- Idempotent consumers via `processed_events` deduplication
- Outbox Pattern for at-least-once event delivery

---

## 8. Communication

- **REST** is the external interface of the Cashback Service.
- **gRPC** is used for internal communication between the Mint Consumer and the Blockchain Adapter.

---

## 9. Trade-offs

- Blockchain is not used for business rules to avoid cost and latency.
- Event-driven flow increases complexity but improves resilience and decoupling.
- Off-chain ledger is not the source of truth for balances; the blockchain is.
- REST simplifies the system entry point.
- gRPC simplifies internal communication but limits direct browser access.
- Redis adds an infrastructure dependency to handle distributed nonce coordination safely.

---

## 10. Intentional Design Decisions

- Separate blockchain interaction into a dedicated adapter.
- Use asynchronous consumers for blockchain operations.
- Persist events before publishing (Outbox Pattern).
- Treat blockchain as the source of truth for token balances.
- Keep the API free of direct blockchain dependencies.
- Use fence tokens in `wallet_nonces` for optimistic concurrency on nonce management.
- Deduplicate consumed events in the Mint Consumer to ensure idempotency.

---

## 11. High-Level Flow Diagram

```
[ Client ]
    |
    | REST
    v
[ Cashback Service ]
    |
    | outbox -> NATS JetStream (cashback.approved)
    v
[ Mint Consumer ]
    |
    | gRPC
    v
[ Blockchain Adapter ]
    |
    | Ethereum JSON-RPC
    v
[ Blockchain / CashbackToken contract ]
```

---

## 12. Purpose of This Repository

This repository serves as:

- a reference implementation
- a validation of Web3, crypto, and blockchain concepts
- a practical backend architecture example
- a foundation for further experimentation and learning
