# Domain Events

## Event Catalog

| Event | Producer | Consumer | Trigger |
|---|---|---|---|
| `purchase.created` | cashback-service-api | cashback-service-api (internal) | Client submits a purchase via REST |
| `cashback.approved` | cashback-service-api (outbox) | mint-consumer | Cashback calculated and approved |
| `token.mint.requested` | mint-consumer | mint-consumer (internal state) | mint-consumer receives `cashback.approved` |
| `token.minted` | mint-consumer | cashback-service-api (optional) | blockchain-adapter confirms successful mint |
| `token.mint.failed` | mint-consumer | mint-consumer (retry logic) | blockchain-adapter returns an error |

---

## Event Payloads

### cashback.approved

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

### token.minted

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

### token.mint.failed

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

## NATS JetStream Configuration

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

## Outbox Pattern

The cashback-service-api uses the Outbox Pattern for reliable event publishing:

1. The business operation and event are written in the **same database transaction**.
2. A separate goroutine polls the `outbox_events` table every 100ms.
3. Events are published to NATS JetStream.
4. Successfully published events are marked as `published`.

This guarantees **at-least-once delivery**: events are never lost even if NATS is temporarily unavailable.

