# E2E Tests

The e2e suite runs the three services against a dedicated isolated stack (Postgres on port 25432, NATS on port 4322, Redis on port 6479). It exercises the full cashback and deposit flows through real HTTP calls and on-chain transactions.


## Running the full suite

```bash
mise run test:e2e
```

The script `scripts/run-e2e.sh` handles the full lifecycle:

1. Starts isolated Docker Compose infrastructure (`docker-compose.e2e.yml`).
2. Waits for Postgres, NATS, and Anvil.
3. Deploys `CashbackToken` to the local Anvil node.
4. Creates NATS streams via `cmd/nats-setup`.
5. Applies DB migrations against the e2e databases.
6. Builds and starts the three services.
7. Waits for the API to be healthy on port 18080.
8. Runs the test suite with `-tags=e2e`.
9. Tears down all infrastructure.


## Anvil EVM node

The e2e infrastructure includes an [Anvil](https://book.getfoundry.sh/anvil/) node (`ghcr.io/foundry-rs/foundry`) running on port `8545`. Anvil is the EVM simulator from the Foundry suite. It provides the same interface as Hardhat Network but ships as a standalone binary with no Node.js dependency.

The node starts with `--chain-id 31337`, matching `ETHEREUM_CHAIN_ID` in `.env`. The deployer account is the standard Foundry test account derived from the `"test test...junk"` mnemonic, which is also the `WALLET_MNEMONIC` in `.env`. After the node is ready, the script deploys `CashbackToken` via Hardhat, producing the deterministic address `0x5FbDB2315678afecb367f032d93F642f64180aa3` already set as `CONTRACT_ADDRESS` in `.env`.


## Test suites

### Cashback suite

`TestCashbackFlow_ShouldIncrementBalanceAfterMint` exercises the mint path. A cashback approval triggers a `CashbackToken.mint()` call and the test verifies the user balance returned by the API increases.

### Deposit suite

`TestDepositFlow_ShouldDetectOnChainTransfer` mints tokens directly on the contract and verifies that `blockchain-adapter` stores the detected `Transfer` event in `detected_deposits` within 45 seconds.

`TestDepositFlow_ShouldCreateDepositReceiptAndCashback` mints tokens, then deposits them to the platform wallet, and verifies that the full chain completes, including detection, deposit receipt creation, and cashback credit.

Both tests skip gracefully if no EVM node is reachable on `ETHEREUM_RPC_URL`.
