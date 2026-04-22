# Web3 Cashback Platform

Backend system that issues cashback as **ERC-20 tokens minted on Ethereum**. Monorepo with three Go services and an event-driven architecture that cleanly separates off-chain business logic from on-chain state.

Users earn cashback on purchases. The cashback value is represented as a crypto token on a blockchain, which enables transparent ownership, on-chain auditability, and a clear boundary between what the backend decides and what the chain guarantees.

> Reference implementation for exploring Web3, blockchain, and crypto concepts through a production-oriented Go backend.

---

## Services

| Service | Role | Protocol |
|---|---|---|
| **cashback-service-api** | Purchases, cashback rules, Outbox Pattern | REST |
| **mint-consumer** | Idempotent event processing, gRPC orchestration | NATS consumer |
| **blockchain-adapter** | Wallet signing, nonce management, ERC-20 minting, deposit monitoring | gRPC server |

## Flows

**Cashback mint**

1. Purchase recorded
2. cashback-service-api calculates the amount and publishes an event via the Outbox Pattern
3. NATS JetStream delivers `cashback.approved` to mint-consumer
4. mint-consumer calls blockchain-adapter via gRPC
5. blockchain-adapter signs and submits `CashbackToken.mint()`
6. CBK tokens credited to the user wallet on-chain

**Deposit detection**

1. blockchain-adapter polls the Ethereum node every 12 seconds for `Transfer` events
2. Each detected transfer is saved to `detected_deposits` (idempotent, de-duplicated by tx hash)
3. blockchain-adapter publishes `deposit.detected` to the `DEPOSIT_EVENTS` NATS stream

---

## Stack

- Go 1.26
- Uber Fx
- Chi
- GORM
- PostgreSQL
- NATS JetStream
- gRPC
- go-ethereum
- Solidity
- Hardhat + Anvil (Foundry)
- OpenTelemetry + Jaeger
- Mise

---

## Prerequisites

- [mise](https://mise.jdx.dev/) task runner
- Go 1.26+
- Docker and Docker Compose
- Node.js 18+ (for Hardhat and contract compilation)

## Quick Start

```bash
# Start infrastructure (PostgreSQL, NATS, Redis)
mise run up

# Create databases and apply migrations (first run only)
mise run db:setup

# Compile the smart contract
cd contracts && npm install && cd ..
mise run compile

# Start a local Ethereum node (keep this terminal open)
mise run evm

# Deploy the contract (second terminal, inside contracts/)
mise run deploy
# Copy the printed contract address to .env

# Generate Go bindings from the ABI
mise run contracts:bindings

# Start all three services (one terminal each)
mise run run:api
mise run run:consumer
mise run run:adapter
```

---

## Documentation

- [architecture.md](docs/architecture.md) - system components and database schema
- [events.md](docs/events.md) - NATS JetStream events and Outbox Pattern
- [guide/](docs/guide/) - development guides (setup, flows, e2e, observability, Sepolia)
- [diagrams/](docs/diagrams/) - Mermaid diagrams
