# Copilot Agent Instructions — cashback-platform

## Automatic Session Bootstrap

On the first user message of every new session, execute these steps before responding.

Skip bootstrap when any of the following is true:
- User message contains "recall" or "/recall" (user is invoking manually — avoid running twice)
- A re-attach prompt is present in context (the `resuming-context` skill handles that case)
- User message is a one-word command or slash command with no codebase context needed (e.g., `/help`, `/clear`)

### Bootstrap Steps

**1. Resolve vault and project:**
```bash
REPO_NAME=$(basename "$(git rev-parse --show-toplevel)" 2>/dev/null || echo "unknown")
echo "VAULT=${COPILOT_VAULT:-NO_VAULT}"
echo "REPO=$REPO_NAME"
```
Vault project folder: `cashback-platform`.
If `COPILOT_VAULT` is not set: skip vault steps and proceed with local files only.

**2. Read active plan:**
```bash
cat .github/plans/web3-cashback-platform-completion/progress.md
```

**3. Check git state:**
```bash
git --no-pager log --oneline -10
git --no-pager status
```

**4. Present summary** (max 20 lines): active plan phase, current branch and recent commits, any uncommitted changes.

---

## Project Overview

Web3 cashback platform that issues cashback as crypto tokens (ERC-20 on Ethereum/Sepolia, SPL on Solana/Devnet).
Event-driven architecture using NATS JetStream and the Outbox Pattern.

### Services

| Service | Entrypoint | Purpose |
|---|---|---|
| `cashback-service-api` | `services/cashback-service-api/cmd/api/` | REST API for purchases and cashback rules |
| `mint-consumer` | `services/mint-consumer/` | Consumes `cashback.approved`, triggers minting via gRPC |
| `blockchain-adapter` | `services/blockchain-adapter/cmd/` | gRPC adapter for Ethereum and Solana interactions |

### Key Shared Packages

| Path | Purpose |
|---|---|
| `kit/clock/` | Mockable time source used by all domains |
| `kit/testsuite/` | Base integration test suite with fixtures and UTC enforcement |
| `proto/` | Protobuf definitions shared across services |
| `cmd/nats-setup/` | NATS stream and consumer provisioning |

### Stack

Go 1.26, PostgreSQL, NATS JetStream, gRPC, Chi, GORM, Uber Fx, Viper, OpenTelemetry + Jaeger, Docker Compose.

---

## Essential Commands

```sh
# Start all infrastructure (DB, NATS, Jaeger, OTel Collector)
docker compose up -d

# Run cashback-service-api
cd services/cashback-service-api && go run ./cmd/api

# Run unit tests (all services)
go test ./...

# Run integration tests (requires Docker infra)
go test -tags=integration ./...

# Run e2e tests
cd test/e2e && go test ./...

# Lint
golangci-lint run ./...
```

---

## Architecture Rules (summary)

Full rules are in `.github/instructions/`. Key points:

- **No cross-domain imports** inside `internal/app/` — domains talk via events or wiring layer only
- **Domain types carry no framework tags** — no `json:` or `gorm:` on domain structs
- **Repository models stay inside the repository** — convert to domain via `ToDomain()` before returning
- **Handlers hold the concrete use case struct** — never define a `UseCase` interface in the handler package
- **Use `dbtx.BaseRepository`** only when the use case owns a `TransactionManager`; use plain `*gorm.DB` otherwise
- **Commits require explicit user authorization** — never commit without a clear "yes" after seeing the commit plan

---

## Codebase Search

Prefer `graphify query` over grep for conceptual searches. Read source files only when editing or when graphify does not have the answer.

| Need | Command |
|---|---|
| Broad context about a concept | `graphify query "concept"` |
| Trace connection between two nodes | `graphify path "A" "B"` |
| Details about a specific node | `graphify explain "NodeName"` |
