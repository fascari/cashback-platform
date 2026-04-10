# Running locally

Step-by-step guide to start all services and execute the full cashback flow on a local machine.
All HTTP requests reference the [Postman collection](../services/cashback-service-api/docs/api/postman_collection.json).

---

## Requirements

- Docker
- Node.js (for the Hardhat local node)
- `.env` at the project root (configure after the contract deploy in step 4)

---

## Step 1: Start infrastructure

```bash
mise run up
```

Confirm all containers are healthy:

```bash
docker ps --format "table {{.Names}}\t{{.Status}}"
```

Expected: `postgres`, `nats`, and `redis` with status `(healthy)`.

---

## Step 2: Apply database migrations

```bash
mise run db:setup
```

Expected: `no change` or migration applied messages per schema (`cashback`, `mint`, `blockchain`).

---

## Step 3: Start the local Hardhat node (terminal 1)

A local EVM node removes the need for Sepolia or any external RPC provider. The default mnemonic `test test test test test test test test test test test junk` is pre-configured and derives the same accounts the adapter uses to sign transactions.

```bash
cd contracts
npx hardhat node
```

Expected: node listening on `http://127.0.0.1:8545` with 20 funded accounts listed.

---

## Step 4: Deploy the contract

In a separate shell (keep the node running in terminal 1):

```bash
cd contracts
npm run deploy:local
```

Copy the contract address printed in the output. Update `.env` at the project root:

```text
ETHEREUM_RPC_URL=http://127.0.0.1:8545
CONTRACT_ADDRESS=0x<address from deploy output>
```

---

## Step 5: Start blockchain-adapter (terminal 2)

```bash
mise run run:adapter
```

Expected log:

```text
gRPC server starting on port 50051
```

---

## Step 6: Start cashback-service-api (terminal 3)

The API must start before the consumer. It creates the NATS JetStream stream `CASHBACK_EVENTS` on startup, which the consumer depends on.

```bash
mise run run:api
```

Expected log:

```text
server started on port 8080
```

---

## Step 7: Start mint-consumer (terminal 4)

```bash
mise run run:consumer
```

Expected log:

```text
cashback consumer started
```

---

## Step 8: Import the Postman collection

Open Postman and import:

```text
services/cashback-service-api/docs/api/postman_collection.json
```

The collection uses a `base_url` variable set to `http://localhost:8080` and auto-captures `user_id`, `purchase_id`, and `cashback_id` from responses via test scripts.

---

## Step 9: Execute the flow

Run requests in this order:

1. **Users > Create User** - creates a user with a wallet address. The `user_id` is saved automatically to the collection variable.

2. **Purchases > Create Purchase** - registers a purchase for the user. The `purchase_id` is saved automatically.

3. **Cashback > Calculate Cashback** - calculates 5% cashback on the purchase amount. The record is saved with `status: approved` and a `cashback.approved` event is published to NATS via the transactional outbox.

---

## Step 10: Observe the mint-consumer (terminal 4)

After step 9.3, the consumer receives the event and logs:

```text
minting cashback cashback_id=X wallet=0x... amount=...
```

It calls the blockchain-adapter via gRPC `MintToken`. On success:

```text
mint completed mint_id=X tx_hash=0x... block=X
```

On a duplicate message (idempotency check):

```text
mint skipped: duplicate event event_id=...
```

---

## Step 11: Verify the result

- **Cashback > Find User Cashback** - confirms the cashback record with `status: approved`.
- **Users > Find User** - confirms the user and wallet address.
- **Users > Get User Balance** - returns the on-chain token balance in wei and the human-readable amount (e.g., `5.00` tokens).

To check the token balance directly via gRPC:

```bash
grpcurl -plaintext \
  -d '{"wallet_address":"0xYourWallet"}' \
  localhost:50051 token.TokenService/GetBalance
```

Expected: balance in wei reflecting the minted tokens.

---

## Retry behavior

If a mint fails (connection error, gas issue), the consumer retries automatically:

- Retry loop runs every 30 seconds
- NATS consumer is configured with `MaxDeliver: 5`
- Failed mints are tracked in the `mint.mint_requests` table with `retry_count` and `error_message`

---

## Teardown

Stop service terminals with `Ctrl+C`, then:

```bash
mise run down
```

---

## What the flow validates in practice

This section describes what actually happens end-to-end when the full flow runs successfully.

**User and purchase creation** go through the REST API and are persisted in PostgreSQL under the `cashback` schema. No side effects occur at this stage.

**Cashback calculation** is the first critical step. When `POST /cashbacks/calculate/{purchase_id}` is called, the use case calculates 5% of the purchase amount and creates a cashback record with `status: approved`. Inside the same database transaction, it writes a `cashback.approved` event to the `cashback.outbox_events` table. This guarantees the event is never lost, even if the process crashes before publishing to NATS.

**Outbox relay** runs as a background goroutine inside cashback-service-api. It polls the outbox table every 5 seconds, picks up pending events, publishes them to the NATS JetStream stream `CASHBACK_EVENTS` on subject `cashback.approved`, and marks each event as `published`. The relay uses a transactional outbox pattern to ensure at-least-once delivery without dual-write risk.

**mint-consumer** receives the message from NATS using a durable pull consumer named `mint-consumer`. It deserializes the payload and calls the `mintcashback` use case, which first checks idempotency: if the `event_id` (a UUID) was already processed, the message is acknowledged and skipped. Otherwise, it creates a `mint_request` record and calls the blockchain-adapter via gRPC.

**blockchain-adapter** receives the `MintToken` gRPC call with the wallet address and token amount in wei. It resolves the nonce for the signing wallet, builds the transaction, signs it with the configured wallet (derived from the mnemonic in `.env`), and submits it to the EVM node. On the local Hardhat network, the transaction confirms immediately. The adapter returns the transaction hash and block number.

**mint-consumer** receives the result, marks the mint as `completed` in the `mint.mint_requests` table, and logs:

```text
mint completed mint_id=X tx_hash=0x... block=X
```

**The token balance** on the recipient wallet increases by the minted amount. This can be verified via `GET /users/{id}/balance`, which calls the blockchain-adapter's `Balance` gRPC method and returns the balance in wei alongside a human-readable token amount.

**Idempotency** is enforced at two levels: the outbox uses the cashback ID as the aggregate key, and the consumer stores each processed `event_id` in a `processed_events` table before creating the mint request. Replaying the same NATS message produces no duplicate mint.

**Retry** is triggered when the gRPC call fails. The mint is marked `failed` with an error code and the consumer retries it every 30 seconds up to 5 attempts. After 5 failures the record stays in `failed` state for inspection.
