# Reference: Web3 Concepts

**Index**: [Ethereum vs Solana](#24-ethereum-vs-solana---quick-comparison) · [Glossary](#25-glossary)

---

## 24. Ethereum vs Solana - Quick Comparison

| Aspect | Ethereum | Solana |
|---|---|---|
| **Consensus** | PoS (Gasper / Casper FFG) | Tower BFT + Proof of History |
| **Finality** | ~12–64 confirmations (~3–13 min) | `finalized` commitment (~13s) |
| **Account model** | EOA + Contract Accounts | All accounts are data accounts |
| **Smart contracts** | Solidity → EVM bytecode | Rust/Anchor → BPF bytecode |
| **Tokens** | ERC-20 (deploy one contract per token) | SPL (one system program, all tokens) |
| **Cryptography** | secp256k1 | ed25519 |
| **Anti-replay** | Nonce (sequential, managed by you) | Recent blockhash (expires ~90s) |
| **Transaction fees** | Gas × gas price (variable, can spike) | ~5,000 lamports (~$0.001, stable) |
| **TPS** | ~15–30 (L1) | ~65,000 theoretical |
| **Reorg risk** | Low (PoS), mitigated by confirmations | Near-zero after `finalized` |
| **Deposit detection** | `FilterLogs` + `Transfer` event | `getSignaturesForAddress` |
| **Go library** | `github.com/ethereum/go-ethereum` | `github.com/gagliardetto/solana-go` |
| **Testnet** | Sepolia | Devnet |
| **Airdrop** | Web faucet | `RequestAirdrop(pubkey, lamports)` |

---

## 25. Glossary

**Ethereum / EVM**

| Term | Meaning |
|---|---|
| **on-chain** | Data or operations that live on the blockchain |
| **off-chain** | Data in a traditional database (PostgreSQL in this project) |
| **tx / txn** | Abbreviation for "transaction" |
| **tx hash** | Unique transaction ID (SHA3 hash of the signed content) |
| **block number** | Sequential number of the block where a tx was confirmed |
| **receipt** | Confirmation that a tx was included in a block |
| **wei** | Smallest ETH unit. 1 ETH = 10^18 wei |
| **deploy** | Publishing a contract to the blockchain |
| **owner** | Address with special permission in the contract |
| **EOA** | Externally Owned Account, wallet controlled by a private key |
| **EVM** | Ethereum Virtual Machine, executes contract bytecode |
| **mempool** | Queue of transactions waiting to be included in a block |
| **finality** | Point after which a transaction can no longer be reversed |
| **testnet** | Test network (e.g., Sepolia). ETH has no real value |
| **mainnet** | Main Ethereum network. ETH has real monetary value |
| **faucet** | Service that distributes testnet ETH for free |
| **ABI** | Application Binary Interface. Defines how to call contract functions |
| **bytecode** | Compiled contract code that runs on the EVM |
| **nonce gap** | Nonce N never used on-chain. Blocks all txs with nonce > N |
| **reorg** | Chain reorganization. An orphaned block reverts its transactions |
| **uncle block** | Ethereum block that was valid but not included in the canonical chain |
| **confirmations** | Number of blocks mined on top of the block containing a tx |
| **safe head** | Ethereum block considered safe after 64 confirmations |
| **finalized checkpoint** | Beacon chain finality. Irreversible after ~2 epochs |

**Solana**

| Term | Meaning |
|---|---|
| **slot** | Solana time unit. A ~400ms window where a validator produces a block |
| **epoch** | ~2–3 day period of ~432,000 slots |
| **lamport** | Smallest SOL unit. 1 SOL = 10^9 lamports |
| **signature** | Solana transaction identifier, ed25519 sig, base58-encoded (88 chars) |
| **account** | Fundamental Solana data unit. Stores bytes and lamport balance |
| **program** | Stateless executable account (Solana's smart contract) |
| **Mint Account** | Account describing an SPL token (supply, decimals, authority) |
| **Token Account** | Account holding a user's balance of a specific SPL token |
| **ATA** | Associated Token Account, deterministic token account address |
| **recent blockhash** | Required in every tx. Expires after ~60-90s |
| **commitment** | Solana finality level: processed / confirmed / finalized |
| **SPL** | Solana Program Library, standard programs including the Token Program |
| **devnet** | Solana development network. Free SOL via `RequestAirdrop` |
| **rent** | Lamports required to keep an account alive (exempt if balance ≥ threshold) |

**Custody and Security**

| Term | Meaning |
|---|---|
| **custody provider** | Service managing private keys (Fireblocks, BitGo) |
| **MPC** | Multi-Party Computation. The key is split into shares, never fully assembled |
| **HSM** | Hardware Security Module, tamper-proof hardware for key storage |
| **cold wallet** | Keys stored offline. Higher security, lower availability |
| **hot wallet** | Keys online and ready to sign. Required for automated operations |
| **fencing token** | Monotonically increasing integer that prevents stale distributed lock writes |
