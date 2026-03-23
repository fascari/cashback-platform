# Web3 Cashback Platform

Backend system that issues cashback as **crypto tokens minted on a blockchain** (Sepolia testnet). Monorepo with three Go services, event-driven architecture, and clean separation between off-chain and on-chain responsibilities.

Users earn cashback when making purchases. Instead of storing it only in a centralized database, the cashback value is represented as an **ERC-20 token on Ethereum**, enabling transparent ownership, on-chain auditability, and clear separation of concerns. Off-chain systems decide what should happen; on-chain systems guarantee the final state.

This repository is a **reference implementation** for exploring Web3, blockchain, and crypto concepts through a production-oriented Go backend.

## Services

| Service                  | Role                                            | Protocol      |
|--------------------------|-------------------------------------------------|---------------|
| **cashback-service-api** | Purchases, cashback rules, Outbox Pattern       | REST          |
| **mint-consumer**        | Idempotent event processing, gRPC orchestration | NATS consumer |
| **blockchain-adapter**   | Wallet signing, nonce management, ERC-20 minting | gRPC server   |

## Stack

Go 1.26 · Uber Fx · Chi · GORM · PostgreSQL · NATS JetStream · gRPC · go-ethereum · Solidity · Hardhat · Mise

## Documentation

- [architecture.md](docs/architecture.md): Components, database schema, communication patterns
- [events.md](docs/events.md): Domain events, NATS JetStream configuration, Outbox Pattern
- [decisions.md](docs/decisions.md): Technology stack, design decisions, trade-offs, fault tolerance
- [diagrams/](docs/diagrams/): PlantUML diagrams
