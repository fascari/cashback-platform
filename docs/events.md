# Domain Events Documentation

## Overview

The system uses an event-driven architecture where domain events define the workflow between services. Events are published to NATS JetStream and consumed asynchronously.

## Event Catalog

### purchase.created

**Description**: A purchase has been registered in the system.

**Producer**: Cashback Service API

**Consumers**: Cashback Service API (internal processing)

**Payload**:
```json
{
  "event_id": "uuid",
  "event_type": "purchase.created",
  "timestamp": "2024-01-15T10:30:00Z",
  "data": {
    "purchase_id": "uuid",
    "user_id": "uuid",
    "amount": 150.00,
    "currency": "USD",
    "merchant_id": "uuid",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

**Trigger**: Client submits a purchase via REST API

**Next Event**: `cashback.approved` (if cashback rules are satisfied)

---

### cashback.approved

**Description**: Cashback has been calculated and approved for a purchase.

**Producer**: Cashback Service API

**Consumers**: Mint Consumer

**Payload**:
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

**Trigger**: Successful cashback calculation after purchase creation

**Next Event**: `token.mint.requested`

---

### token.mint.requested

**Description**: A request to mint tokens has been issued.

**Producer**: Mint Consumer

**Consumers**: Mint Consumer (internal state tracking)

**Payload**:
```json
{
  "event_id": "uuid",
  "event_type": "token.mint.requested",
  "timestamp": "2024-01-15T10:30:02Z",
  "data": {
    "mint_request_id": "uuid",
    "cashback_id": "uuid",
    "user_id": "uuid",
    "wallet_address": "0x...",
    "token_amount": "1500000000000000000",
    "idempotency_key": "uuid"
  }
}
```

**Trigger**: Mint Consumer receives `cashback.approved` event

**Next Event**: `token.minted` or `token.mint.failed`

---

### token.minted

**Description**: Tokens were successfully minted on-chain.

**Producer**: Mint Consumer (after successful gRPC call to Blockchain Adapter)

**Consumers**: Cashback Service API (optional, for ledger update)

**Payload**:
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

**Trigger**: Blockchain Adapter confirms successful mint

**Next Event**: None (terminal event)

---

### token.mint.failed

**Description**: Token minting failed and may be retried.

**Producer**: Mint Consumer (after failed gRPC call to Blockchain Adapter)

**Consumers**: Mint Consumer (for retry logic)

**Payload**:
```json
{
  "event_id": "uuid",
  "event_type": "token.mint.failed",
  "timestamp": "2024-01-15T10:30:05Z",
  "data": {
    "mint_request_id": "uuid",
    "cashback_id": "uuid",
    "user_id": "uuid",
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

**Trigger**: Blockchain Adapter returns an error

**Next Event**: `token.mint.requested` (retry) or dead letter (max retries exceeded)

---

## Event Flow Diagram

```
┌──────────────────────────────────────────────────────────────────────────┐
│                           EVENT FLOW                                      │
└──────────────────────────────────────────────────────────────────────────┘

Client Request
      │
      ▼
┌─────────────────┐
│ purchase.created │ ──────────────────────────────────────────────┐
└─────────────────┘                                                │
      │                                                            │
      │ (cashback calculation)                                     │
      ▼                                                            │
┌─────────────────┐                                                │
│ cashback.approved│ ─────────────────────────────────┐            │
└─────────────────┘                                   │            │
      │                                               │            │
      │ (consumed by Mint Consumer)                   │            │
      ▼                                               │            │
┌─────────────────────┐                               │            │
│ token.mint.requested │                              │            │
└─────────────────────┘                               │            │
      │                                               │            │
      │ (gRPC call to Blockchain Adapter)             │            │
      ▼                                               │            │
┌─────────────────────────────────────┐               │            │
│                                     │               │            │
│  ┌─────────────┐  ┌───────────────┐ │               │            │
│  │ token.minted │  │token.mint.failed│ │               │            │
│  └─────────────┘  └───────────────┘ │               │            │
│       │                   │         │               │            │
│       │                   │         │               │            │
│       ▼                   ▼         │               │            │
│    SUCCESS             RETRY?       │               │            │
│                          │          │               │            │
│                    ┌─────┴─────┐    │               │            │
│                    │           │    │               │            │
│                   YES          NO   │               │            │
│                    │           │    │               │            │
│                    ▼           ▼    │               │            │
│               (retry)    (dead letter)              │            │
│                                     │               │            │
└─────────────────────────────────────┘               │            │
                                                      │            │
                                                      │            │
Producer Legend:                                      │            │
─────────────────                                     │            │
• purchase.created    → Cashback Service API          │            │
• cashback.approved   → Cashback Service API ─────────┘            │
• token.mint.requested → Mint Consumer                             │
• token.minted        → Mint Consumer                              │
• token.mint.failed   → Mint Consumer ─────────────────────────────┘
```

## NATS JetStream Configuration

### Streams

```
Stream: CASHBACK_EVENTS
├── Subjects: cashback.>
├── Retention: Limits
├── MaxAge: 7 days
├── Storage: File
└── Replicas: 1 (for dev), 3 (for prod)

Stream: TOKEN_EVENTS
├── Subjects: token.>
├── Retention: Limits
├── MaxAge: 7 days
├── Storage: File
└── Replicas: 1 (for dev), 3 (for prod)

Stream: PURCHASE_EVENTS
├── Subjects: purchase.>
├── Retention: Limits
├── MaxAge: 7 days
├── Storage: File
└── Replicas: 1 (for dev), 3 (for prod)
```

### Consumers

```
Consumer: mint-consumer
├── Stream: CASHBACK_EVENTS
├── FilterSubject: cashback.approved
├── DeliverPolicy: All
├── AckPolicy: Explicit
├── MaxDeliver: 5
└── AckWait: 30s

Consumer: cashback-service-token-updates
├── Stream: TOKEN_EVENTS
├── FilterSubject: token.minted
├── DeliverPolicy: All
├── AckPolicy: Explicit
├── MaxDeliver: 5
└── AckWait: 30s
```

## Idempotency

All event consumers must be idempotent. This is achieved through:

1. **Event ID tracking**: Each event has a unique `event_id`
2. **Processed events table**: Consumers track processed event IDs
3. **Idempotency keys**: Business operations include idempotency keys

### Idempotency Check Flow

```
1. Receive event
2. Check if event_id exists in processed_events table
3. If exists → Acknowledge and skip
4. If not exists → Process event
5. On success → Insert event_id into processed_events
6. Acknowledge event
```

## Outbox Pattern Implementation

The Cashback Service uses the Outbox Pattern:

```
┌────────────────────────────────────────────────────────────────────┐
│                    OUTBOX PATTERN FLOW                              │
└────────────────────────────────────────────────────────────────────┘

1. Business Transaction
   ┌─────────────────────────────────────────┐
   │ BEGIN TRANSACTION                        │
   │   INSERT INTO purchases (...)            │
   │   INSERT INTO cashback_ledger (...)      │
   │   INSERT INTO outbox_events (...)        │
   │ COMMIT                                   │
   └─────────────────────────────────────────┘

2. Outbox Publisher (separate goroutine)
   ┌─────────────────────────────────────────┐
   │ LOOP every 100ms:                        │
   │   SELECT * FROM outbox_events            │
   │     WHERE status = 'pending'             │
   │     ORDER BY created_at                  │
   │     LIMIT 100                            │
   │                                          │
   │   FOR each event:                        │
   │     Publish to NATS JetStream            │
   │     UPDATE outbox_events                 │
   │       SET status = 'published'           │
   │       WHERE id = ?                       │
   └─────────────────────────────────────────┘
```

## Retry Strategy

Failed events follow exponential backoff:

| Retry | Delay |
|-------|-------|
| 1 | 1s |
| 2 | 2s |
| 3 | 4s |
| 4 | 8s |
| 5 | 16s |

After 5 retries, events are moved to a dead letter queue for manual inspection.

