# Web3 Cashback Platform – Project Specification

## 1. Project Overview

This document specifies a backend system that issues cashback as **crypto tokens minted on a blockchain**.

The system integrates traditional backend components with **Web3 concepts**, using an **event-driven architecture** and a clear separation between **off-chain** and **on-chain** responsibilities.

This repository is intended to **validate and explore business, Web3, crypto, and blockchain concepts** through a pragmatic and production-oriented design.

---

## 2. Core Business Idea

Users earn cashback when making purchases.

Instead of storing cashback only in a centralized database, the cashback value is represented as a **crypto token issued on a blockchain**, enabling:

- transparent ownership
- auditability
- clear separation of responsibilities

---

## 3. Web3, Blockchain, and Crypto Concepts

### Blockchain
A **distributed and immutable ledger** used as the **source of truth** for token ownership.

### Crypto Token
A **digital asset** issued on a blockchain, representing the cashback value earned by users.

### Web3
An application model where **users directly own digital assets**, while backend systems act as coordinators rather than custodians.

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

> **Off-chain systems decide what should happen.  
> On-chain systems guarantee the final and auditable state.**

---

## 5. System Architecture Components

- **Cashback Service**  
  Core service responsible for purchases, cashback rules, and off-chain persistence.

- **Mint Consumer**  
  Asynchronous component that processes events and triggers token minting.

- **Blockchain Adapter**  
  Isolated service that abstracts blockchain interactions via gRPC.

- **Event Broker**  
  NATS JetStream used for reliable and durable event delivery.

- **Database**  
  PostgreSQL used for off-chain state (**one database per service**).

---

## 6. Database Tables (Conceptual, Per Service)

Each service manages its own database schema.

Conceptual tables include:

- **users**  
  Users and associated wallet addresses.

- **purchases**  
  Purchase records.

- **cashback_ledger**  
  Off-chain representation of generated cashback.

- **blockchain_transactions**  
  Status of on-chain mint operations.

- **outbox_events**  
  Pending events to be published.

---

## 7. Domain Events

The following domain events define the system workflow:

- **purchase.created**  
  A purchase has been registered.

- **cashback.approved**  
  Cashback amount has been calculated and approved.

- **token.mint.requested**  
  A request to mint tokens has been issued.

- **token.minted**  
  Tokens were successfully minted on-chain.

- **token.mint.failed**  
  Token minting failed and may be retried.

### Fault Tolerance

- Retry
- Replay
- Fallback
- Idempotent consumers

---

## 8. Communication

- **REST** is used as the external interface of the Cashback Service.
- **gRPC** is used for internal service-to-service communication.

---

## 9. Trade-offs

- Blockchain is not used for business rules to avoid cost and latency.
- Event-driven flow increases complexity but improves resilience and decoupling.
- Off-chain ledger is a cache and not the source of truth.
- REST simplifies the system entry point.
- gRPC simplifies internal communication but limits direct browser access.

---

## 10. Intentional Design Decisions

- Separate blockchain interaction into a dedicated adapter.
- Use asynchronous consumers for blockchain operations.
- Persist events before publishing (Outbox Pattern).
- Treat blockchain as the source of truth for balances.
- Keep the API free of direct blockchain dependencies.

---

## 11. High-Level Flow Diagram

```
[ Client ]
    |
    v
[ Cashback Service ]
    |
    | (outbox event)
    v
[ NATS JetStream ]
    |
    v
[ Mint Consumer ]
    |
    v
[ Blockchain Adapter ]
    |
    v
[ Blockchain ]
```

---

## 12. Purpose of This Repository

This repository serves as:

- a **reference implementation**
- a **validation of Web3, crypto, and blockchain concepts**
- a **practical backend architecture example**
- a **foundation for further experimentation and learning**
