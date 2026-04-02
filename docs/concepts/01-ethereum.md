# Ethereum — Web3 Concepts

**Index**: [Blockchain](#1-blockchain) · [Ethereum](#2-ethereum) · [Wallet](#3-wallet-and-key-management) · [Gas](#4-gas) · [Nonce](#5-nonce) · [TX Lifecycle 🎯](#6-ethereum-transaction-lifecycle-) · [Smart Contract](#7-smart-contract) · [Solidity](#8-solidity) · [ERC-20](#9-erc-20) · [ABI & abigen](#10-abi-and-abigen) · [Sepolia](#11-sepolia-testnet) · [RPC Providers](#12-rpc-providers-infura--alchemy) · [Hardhat](#13-hardhat)

---

## 1. Blockchain

A blockchain is a distributed, immutable database. Each "block" holds a batch of transactions.
Once confirmed, a block cannot be changed, unlike a SQL database where `UPDATE` is always possible.

**Go analogy**: a distributed `append-only log` where each entry is cryptographically linked to
the previous one. No `DELETE`, no `UPDATE`, only `INSERT`.

**Impact on code**: on-chain operations are asynchronous and irreversible. You submit a transaction
and wait for confirmation. There is no rollback. That is why the project tracks every mint with an
`idempotency_key` in `blockchain_transactions`.

**Where it appears**: `blockchain-adapter` handles every mint operation through the chain.

---

## 2. Ethereum

Ethereum is a programmable blockchain. Beyond transferring value (ETH), you can publish code
that runs on the network, called **smart contracts**. It is a distributed computer whose
"processor" is the EVM (Ethereum Virtual Machine).

**Where it appears**: the `blockchain-adapter` mints ERC-20 tokens on Ethereum via `go-ethereum`.
Every mint goes through the EVM on the Sepolia testnet.

---

## 3. Wallet and Key Management

A wallet does not store tokens. It stores a **private key**. The private key is a 256-bit number
that mathematically derives a **public address** (e.g., `0xAbCd...`).

```text
mnemonic (12 words) -> seed -> private key -> public key -> address (0x...)
```

| Component | Description |
|---|---|
| Private key | Secret. Signs transactions. Whoever holds this has full control. |
| Public address | The "account number". Can be shared freely. |
| Mnemonic | 12 or 24 words that regenerate the private key. The human-readable backup. |

**HD Wallet (BIP-32/BIP-39/BIP-44)**: a single mnemonic can deterministically derive multiple
addresses via a derivation path. The path encodes coin type, account, and index:

```
Ethereum: m/44'/60'/0'/0/0   (coin_type=60)
Solana:   m/44'/501'/0'/0'   (coin_type=501)
```

This allows managing many wallets from one seed — critical for exchange wallets where each
user gets a unique deposit address derived from the same master key.

**Cryptographic algorithms by chain**:

| Chain | Algorithm | Key size | Address format |
|---|---|---|---|
| Ethereum | secp256k1 | 32 bytes | 0x + 20 bytes hex |
| Bitcoin | secp256k1 | 32 bytes | Base58 |
| Solana | ed25519 | 32 bytes secret + 32 bytes public | Base58 (44 chars) |
| XRP | secp256k1 or ed25519 | 32 bytes | Base58 with checksum |

**Production key management**: in production, private keys are never stored in environment
variables or code. They live in **custody providers** (Fireblocks, BitGo) or HSMs
(Hardware Security Modules). See [Custody Providers](03-production-patterns.md#21-custody-providers--fireblocks-and-bitgo-).

**Where it appears**: `internal/infra/ethereum/wallet.go` in `blockchain-adapter`.
The service derives the wallet from a development mnemonic to sign mint transactions.

---

## 4. Gas

Every EVM operation costs **gas**, a unit of computational work. You pay in ETH.
The gas price varies with network demand.

```text
total cost = gasUsed x gasPrice
```

**EIP-1559**: since EIP-1559, the gas price is split into `baseFee` (burned by the network)
and `maxPriorityFee` (tip for the validator). `SuggestGasPrice()` from go-ethereum returns
a value appropriate for current network conditions.

**Where it appears**: `blockchain-adapter` calls `ethClient.SuggestGasPrice()` before submitting
each transaction. The estimated gas for a `mint()` call is low (~50k gas units).

---

## 5. Nonce

Nonce (number used once) is a sequential counter per address. The first transaction from a wallet
has nonce=0, the second nonce=1, and so on.

**Why it exists**: prevents replay attacks (resending the same transaction) and enforces ordering.
The network rejects transactions with a nonce already used or one that creates a gap in the sequence.

**Problem in concurrent systems**: two goroutines trying to mint simultaneously may both read
nonce=5. One will fail with "nonce too low". The solution in this project is a distributed lock
with a fencing token via Redis, described in the implementation plan, Phase 1.3.

**Nonce gap**: if `GetAndIncrement` increments the nonce in the DB but `SendTransaction` fails,
nonce N is never used on-chain. The network rejects nonce N+1 because N was never confirmed.
Fix: call `SyncFromChain()` after any `SendTransaction` failure.

**Where it appears**: `internal/repository/nonce/repository.go` in `blockchain-adapter`.

---

## 6. Ethereum Transaction Lifecycle 🎯

A transaction is a signed instruction sent to the network. Main fields:

| Field | Description |
|---|---|
| `to` | Contract address or recipient |
| `value` | ETH to transfer (0 for contract calls) |
| `data` | ABI-encoded function call payload |
| `nonce` | Sequential sender counter |
| `gasLimit` | Maximum gas willing to spend |
| `gasPrice` / `maxFeePerGas` | Price per gas unit (EIP-1559) |
| `signature` | Cryptographic proof of who signed (derived from the private key) |

**Full lifecycle with production states**:

```text
pending -> submitted -> confirming (waiting N blocks) -> confirmed
                     \> failed (reverted / out of gas)
                     \> dropped (nonce too low / replaced)
                     \> reorged (block containing tx was orphaned)
```

This project tracks this lifecycle in `blockchain_transactions.status`.

**Transaction hash**: SHA3 hash of the signed transaction content. Deterministic — given the
same nonce and payload, the hash is always the same. This allows recovering a lost transaction
without re-signing.

**Why `confirming` state matters**: a transaction included in a block is not final. The block
itself can be orphaned during a chain reorganization (see [Chain Reorganization](03-production-patterns.md#18-chain-reorganization-)).
Production systems wait for a minimum number of confirmations before crediting a user's account.

| Use case | Safe confirmation count |
|---|---|
| Low-value transactions | 3–6 confirmations |
| Exchange deposits / payments | 12–64 confirmations |
| Ethereum "safe head" | 64 confirmations (~13 min) |
| Ethereum "finalized" (PoS) | ~2 epochs = 64 slots (~13 min) |

**Where it appears**: `internal/domain/transaction.go` in `blockchain-adapter`.

---

## 7. Smart Contract

A smart contract is code published on the blockchain. Once published, it has a permanent address
and its code cannot be changed (immutability). Anyone can call its public functions.

**Go analogy**: a struct with methods that lives on the blockchain, has persistent on-chain state,
and every call is a paid transaction.

```solidity
contract CashbackToken {
    mapping(address => uint256) balances; // on-chain state

    function mint(address to, uint256 amount) external {
        balances[to] += amount; // mutates on-chain state — irreversible
    }
}
```

**Where it appears**: `contracts/CashbackToken.sol`, ~20 lines. The goal is not to master
Solidity, but to understand the contract well enough to integrate it from Go.

---

## 8. Solidity

Programming language for EVM smart contracts. Syntax is close to C++/Java.
Relevant elements for this project:

| Element | Go equivalent |
|---|---|
| `address` | `string` (20-byte hex address) |
| `uint256` | `*big.Int` |
| `mapping(address => uint256)` | `map[string]*big.Int` |
| `external` | public function callable from outside the contract |
| `onlyOwner` | access control modifier |

**Where it appears**: `contracts/CashbackToken.sol`, ~20 lines. The goal is not to master
Solidity, but to understand the contract well enough to integrate it from Go.

---

## 9. ERC-20

ERC-20 is an interface standard for fungible tokens on Ethereum. It defines the functions
every compliant token must implement:

| Function | Description |
|---|---|
| `balanceOf(address)` | Returns the balance of an address |
| `transfer(to, amount)` | Transfers tokens from caller to `to` |
| `totalSupply()` | Total tokens in circulation |
| `mint(to, amount)` | Creates new tokens (common extension, not required by the standard) |

**Fungible tokens**: every unit is identical to every other, like cents of a currency.
Different from NFTs (ERC-721), where each token is unique.

**Where it appears**: `CashbackToken` is an ERC-20 where only the owner (the service wallet)
can call `mint`. When a cashback is approved, tokens are minted to the user's wallet.

---

## 10. ABI and abigen

**ABI** (Application Binary Interface): specification of how to call a contract's functions.
Equivalent to a `.proto` file for gRPC, it defines types and signatures in JSON.

```json
{
  "name": "mint",
  "inputs": [
    { "name": "to",     "type": "address" },
    { "name": "amount", "type": "uint256" }
  ]
}
```

**abigen**: go-ethereum tool that reads an ABI and generates typed Go code. Equivalent to
`protoc` for Solidity contracts.

```text
CashbackToken.sol
  -> npx hardhat compile -> ABI + bytecode
  -> abigen             -> internal/contracts/cashbacktoken.go
```

**Where it appears**: `mise run contracts:bindings` generates
`services/blockchain-adapter/internal/contracts/cashbacktoken.go`.
The generated code exposes `cashbackToken.Mint(opts, to, amount)` with full type safety.

---

## 11. Sepolia Testnet

Sepolia is the official Ethereum testnet. It is a parallel network identical to mainnet,
but with ETH that has no real value, obtained for free from faucets. Used for development and testing.

**Why Sepolia and not mainnet**: every mainnet transaction costs real ETH. On Sepolia,
ETH has no monetary value, so transactions carry no real cost.

**Faucet**: `sepoliafaucet.com` distributes SepoliaETH for free.
~0.1 SepoliaETH is enough to deploy the contract and run mints during development.

**Where it appears**: `ETHEREUM_RPC_URL=https://sepolia.infura.io/v3/{API_KEY}` in `.env`.

---

## 12. RPC Providers (Infura / Alchemy)

To interact with the blockchain, you need an **Ethereum node**. Running a node locally
requires ~1TB of disk and days of sync time. Infura and Alchemy provide node access via API.

```text
blockchain-adapter -> ethclient.Dial(RPC_URL) -> Infura/Alchemy -> Ethereum node -> Sepolia
```

**Where it appears**: `internal/infra/ethereum/client.go` in `blockchain-adapter`.
`ethclient.Dial(cfg.EthereumRPCURL)` connects to the configured provider.

---

## 13. Hardhat

JavaScript/TypeScript framework for smart contract development. In this project it is used
only to compile Solidity and deploy the contract.

```text
mise run contracts:compile        -> npx hardhat compile -> ABI + bytecode
mise run contracts:bindings       -> abigen              -> Go bindings
mise run contracts:deploy-local   -> Hardhat local network
mise run contracts:deploy-sepolia -> Sepolia testnet
```

Go cannot compile Solidity natively. Hardhat is the intermediate compiler.

**Where it appears**: `contracts/` at the repository root.
