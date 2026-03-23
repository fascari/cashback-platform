# Transaction Lifecycle

State machine from MintRequest creation to on-chain confirmation.

```mermaid
stateDiagram-v2
    [*] --> pending : MintRequest created
    pending --> processing : gRPC MintToken called
    processing --> submitted : tx signed and sent
    submitted --> confirmed : block mined, receipt received
    submitted --> failed : RPC error or tx reverted
    failed --> processing : RetryFailedMints, exponential backoff
    confirmed --> [*] : token exists on-chain
    failed --> [*] : max_retries exceeded
```

| State | Fields stored |
|---|---|
| `pending` | `idempotency_key`, `wallet_address`, `token_amount` |
| `submitted` | + `transaction_hash`, `nonce`, `gas_price` |
| `confirmed` | + `block_number`, `gas_used`, `confirmed_at` |
| `failed` | + `error_code`, `error_message` |

