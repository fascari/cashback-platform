# Technology Stack, Design Decisions, and Trade-offs

## Technology Stack

| Concern | Technology | Rationale |
|---|---|---|
| Language | Go 1.26 | Performance, concurrency, simplicity |
| DI | Uber Fx | Lifecycle management, testability |
| Config | Viper | Multi-source config, env var support |
| HTTP | Chi | Lightweight, stdlib compatible |
| ORM | GORM | Productivity, migration support |
| Messaging | NATS JetStream | Persistence, replay, at-least-once delivery |
| RPC | gRPC + protobuf | Type safety, performance, streaming |
| Database | PostgreSQL | ACID, reliability, JSONB support |
| Blockchain (Ethereum) | go-ethereum | Official Go Ethereum client (ethclient, abigen) |
| Blockchain (Solana) | github.com/gagliardetto/solana-go v1.10.0 | Full-featured Solana Go SDK (SPL Token, devnet) |
| Smart Contract | Solidity + OpenZeppelin | ERC-20 standard, battle-tested |
| Contract Tooling | Hardhat | Compilation, local network, deploy scripts |
| Distributed Lock | Redis | Nonce serialization with fencing token (multi-replica safe) |
| Task Runner | Mise | Polyglot task runner |
| Testing | Mockery v2 | Mock generation from interfaces |

---

## Design Decisions

| Decision | Rationale |
|---|---|
| Blockchain interaction isolated in dedicated adapter | Single point of change; other services stay blockchain-agnostic |
| `ChainClient` interface + `ChainRegistry` in blockchain-adapter | Multi-chain support (Ethereum + Solana) without changing upper layers; routed by `chain_id` |
| Async consumers for blockchain operations | On-chain calls are slow and unreliable; async decouples availability |
| Outbox Pattern for event publishing | Guarantees at-least-once delivery without distributed transactions |
| Blockchain as source of truth for balances | Off-chain ledger is a cache; on-chain state is the authority |
| Deposit monitoring via polling (not webhooks) | Simpler operational model; no public endpoint exposure; compatible with any RPC node |
| Confirmation-before-credit policy | Prevents acting on unconfirmed transactions; configurable per chain (ETH: N blocks, Solana: `finalized` commitment) |
| Chain reorg handling in blockchain-adapter | Ethereum reorgs are a production reality; transactions must transition to `reorged` and be re-submitted |
| REST for external API | Simple, widely supported, no browser limitations |
| gRPC for internal communication | Type-safe, fast, streaming support |
| One database per service | No shared state; services communicate only via events or gRPC |

---

## Trade-offs

- Blockchain is not used for business rules: avoids cost and latency.
- Event-driven flow increases complexity but improves resilience and decoupling.
- Off-chain ledger is a read cache, not the source of truth.
- gRPC simplifies internal communication but limits direct browser access.
- `ChainClient` abstraction adds an indirection layer — worth the complexity for multi-chain support.
- Polling-based deposit monitor adds infrastructure overhead vs event subscriptions, but avoids WebSocket fragility.

---

## Fault Tolerance

### Idempotency

- All consumers are idempotent via `event_id` tracking in `processed_events`.
- `idempotency_key` in `blockchain_transactions` prevents duplicate on-chain mints.
- Events can be safely replayed.

### Nonce Serialization (Ethereum)

Ethereum requires strictly sequential nonces per wallet. Using `SELECT FOR UPDATE` fails under multiple service replicas because PostgreSQL row locks do not span across connections from different pods.

**Solution**: Redis distributed lock with a fencing token.

1. Acquire lock with TTL (e.g., 10s) — returns a fencing token (monotonically increasing value).
2. Read current nonce from `wallet_nonces`.
3. Submit transaction with that nonce.
4. Increment nonce in `wallet_nonces` (only if fencing token is still valid).
5. Release lock.

If the service crashes after step 3 but before step 4, the next lock acquisition re-syncs the nonce from the chain via `eth_getTransactionCount`.

### Retry Strategy

Failed operations follow exponential backoff:

| Retry | Delay |
|---|---|
| 1 | 1s |
| 2 | 2s |
| 3 | 4s |
| 4 | 8s |
| 5 | 16s |

After 5 retries, events are moved to a dead-letter queue for manual inspection.

### Fallback

- Temporary failures do not break the system flow.
- Events remain in queues until successfully processed.
- Circuit breaker pattern for blockchain RPC calls.
