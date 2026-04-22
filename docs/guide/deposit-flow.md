# Deposit Detection Flow

The blockchain-adapter polls the Ethereum node for ERC-20 Transfer events. When a transfer targets the platform wallet, it is recorded as a deposit and the cashback-service-api credits the user.

This flow runs in the opposite direction of the cashback flow. Instead of the platform writing to the blockchain, it reads from it and reacts.

See [setup.md](setup.md) to start the stack before running this flow.


## How it works

1. A `Transfer` event is emitted on the `CashbackToken` contract (from any source: mint, transfer, or external wallet).
2. `blockchain-adapter` polls the Ethereum node every 12 seconds for new Transfer events in the latest block range.
3. Each detected event is saved to `blockchain.detected_deposits` with its transaction hash, wallet address, token amount, and block number.
4. A `deposit.detected` event is published to the `DEPOSIT_EVENTS` NATS stream.
5. `cashback-service-api` deposit consumer receives the event, looks up the user by wallet address, creates a `deposit_receipt`, and credits 1% cashback.

The monitor resumes from the last stored block number on restart, so no events are lost after a service restart.

---


## Testing manually

With the stack running and the contract deployed, trigger a Transfer event by calling `mint()` from the deployer account:

```bash
cast send \
  --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
  --rpc-url http://127.0.0.1:8545 \
  <CONTRACT_ADDRESS> \
  "mint(address,uint256)" \
  0x70997970C51812dc3A010C7d01b50e0d17dc79C8 \
  1000000000000000000
```

Replace `<CONTRACT_ADDRESS>` with the address printed during `mise run deploy`.

The recipient `0x70997970C51812dc3A010C7d01b50e0d17dc79C8` is the well-known Anvil account #1. It must have a user record in the cashback database with that wallet address for the deposit receipt to be created.

The blockchain-adapter polls every 12 seconds. Within one cycle, a row appears in `blockchain.detected_deposits` and a `deposit.detected` event is published to the `DEPOSIT_EVENTS` NATS stream.


### Verify detection in blockchain-adapter

```bash
psql "$POSTGRES_DSN_BLOCKCHAIN" \
  -c "SELECT transaction_hash, wallet_address, token_amount, block_number, status FROM blockchain.detected_deposits ORDER BY id DESC LIMIT 5;"
```


### Verify deposit receipt in cashback-service-api

```bash
psql "$POSTGRES_DSN_CASHBACK" \
  -c "SELECT tx_hash, from_address, amount, block_number FROM cashback.deposit_receipts ORDER BY id DESC LIMIT 5;"
```


### Verify cashback credited from the deposit

```bash
psql "$POSTGRES_DSN_CASHBACK" \
  -c "SELECT id, user_id, amount, status, source FROM cashback.cashback_ledger WHERE source = 'deposit' ORDER BY id DESC LIMIT 5;"
```

---


## Difference from the cashback flow

| | Cashback flow | Deposit flow |
|---|---|---|
| Trigger | REST API call | On-chain Transfer event |
| Blockchain direction | Write (mint tokens) | Read (detect Transfer) |
| Adapter role | Executes mint via gRPC | Polls for events and publishes to NATS |
| Cashback rate | 5% of purchase amount | 1% of deposited token amount |

The cashback flow validates that the system can write to the blockchain. The deposit flow validates that it can read and react to on-chain events. Together they exercise the full bidirectional interaction.
