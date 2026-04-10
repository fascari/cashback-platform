# Solana — Web3 Concepts

**Index**: [Solana](#14-solana) · [SPL Token Program](#15-spl-token-program) · [Transaction & Commitment Levels](#16-solana-transaction-and-commitment-levels) · [Devnet](#17-solana-devnet)

---

## 14. Solana

Solana is a high-throughput blockchain using **Proof of History (PoH)** combined with
**Tower BFT** consensus. Key differences from Ethereum:

**Account Model**: everything on Solana is an *Account*. Accounts store data (raw bytes) and
a balance in lamports. Programs (smart contracts) are stateless accounts with `executable=true`.
State lives in separate *data accounts* that the program controls — not inside the program itself.

```
Ethereum:
  ERC-20 Contract (0xABC...)
  └── balances[0xUser] = 100   ← state inside the contract

Solana:
  Token Program (system, pre-deployed) ← stateless program
  Mint Account (7xKX...)               ← describes the token
  Token Account (3aZZ...)              ← holds a user's balance
```

**Performance**: ~65,000 theoretical TPS, ~400ms block time (slots), transaction fees of
~5,000 lamports (~$0.001). No gas — fees are flat and predictable.

**Finality**: Solana uses three commitment levels instead of block confirmations (see Section 16).
Once `finalized`, transactions are mathematically irreversible.

**Where it appears**: `internal/infra/solana/` in `blockchain-adapter` (Etapa 1B).

---

## 15. SPL Token Program

**SPL (Solana Program Library) Token** is the Solana equivalent of ERC-20 — but instead of
deploying a new contract per token, all tokens share one system-level program
(`TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA`).

**Three account types in SPL**:

| Account | Purpose | Analogy |
|---|---|---|
| **Mint Account** | Describes the token: decimals, total supply, mint authority | ERC-20 contract address |
| **Token Account** | Holds a specific user's balance for one token | `balances[user]` inside ERC-20 |
| **Associated Token Account (ATA)** | Token Account at a deterministic address derived from (owner, mint) | Computed mapping |

**ATA (Associated Token Account)**: the address is fully deterministic:

```
ATA = findProgramAddress([owner_pubkey, TOKEN_PROGRAM_ID, mint_pubkey], ASSOCIATED_TOKEN_PROGRAM_ID)
```

Given a wallet address and a mint address, you always know the ATA — no storage needed.
This is a major design improvement over ERC-20 where the contract manages the mapping.

**Before minting**: if the recipient's ATA does not exist, it must be created first.
The `CreateAssociatedTokenAccountIdempotent` instruction does this atomically — it is a no-op
if the ATA already exists, making it safe to include in every mint transaction.

**MintTo flow**:
```
Transaction = [
  Instruction 1: CreateAssociatedTokenAccountIdempotent (no-op if exists)
  Instruction 2: MintTo(mint_account, destination_ata, amount)
]
```

Both instructions execute atomically in one transaction.

**Where it appears**: `internal/infra/solana/spl_token.go` in `blockchain-adapter` (Etapa 1B).

---

## 16. Solana Transaction and Commitment Levels

**No nonce — Recent Blockhash instead**: Solana transactions do not use a sequential nonce.
Instead, every transaction must include a `recentBlockhash` — the hash of a recent slot.
A blockhash is valid for ~60–90 seconds. After that, the transaction is rejected.

```go
hash, _ := client.GetRecentBlockhash(ctx, commitment.Confirmed)
tx := solana.NewTransaction(instructions, hash.Value.Blockhash, feePayer)
```

This eliminates nonce management complexity, but introduces a different risk: if a transaction
is not submitted or confirmed within the validity window, it expires and must be re-built with
a fresh blockhash.

**Commitment Levels** (Solana's equivalent of block confirmations):

| Level | Meaning | Safety | Latency |
|---|---|---|---|
| `processed` | Included in a slot, not yet voted on | Unsafe — can be reverted | ~400ms |
| `confirmed` | Supermajority of validators voted for the slot | 99%+ safe | ~1–2s |
| `finalized` | Slot has reached max lockout (~32 votes) | **Irreversible** | ~13s |

**Practical use**:
- Use `confirmed` for UX (showing a pending state to the user)
- Use `finalized` before crediting accounts or releasing funds

**Signature vs Hash**: Solana uses `signature` instead of `tx_hash`. A signature is
the ed25519 signature of the transaction's message — 64 bytes, base58-encoded (88 chars).
This is the canonical transaction identifier in `getTransaction` and `getSignatureStatuses`.

**Where it appears**: `internal/infra/solana/finality_watcher.go` (Etapa 1B).

---

## 17. Solana Devnet

Solana's development testnet, equivalent to Ethereum's Sepolia.

| Property | Sepolia (Ethereum) | Devnet (Solana) |
|---|---|---|
| RPC URL | `https://sepolia.infura.io/v3/{key}` | `https://api.devnet.solana.com` |
| Native token | SepoliaETH | Devnet SOL |
| Getting tokens | Web faucet | Programmatic airdrop |
| Explorer | `sepolia.etherscan.io` | `explorer.solana.com/?cluster=devnet` |

**Programmatic airdrop** (Solana advantage):

```go
sig, _ := client.RequestAirdrop(ctx, pubkey, 1_000_000_000) // 1 SOL
```

No manual web interaction needed — useful for automated test setup.
