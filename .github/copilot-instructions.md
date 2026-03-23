# GitHub Copilot Instructions — cashback-platform

This file is the entry point for AI agents working in this repo. It provides project context and points to the canonical sources for rules, conventions, and workflows. Do not duplicate content that already exists in `.github/ai/rules/` or `.github/ai/skills/`.

---

## Project Basics

- **Language**: Go (multiple services with independent modules)
- **Services**:
  - `cashback-service-api`: REST API for cashback operations
  - `blockchain-adapter`: gRPC service for blockchain interactions
  - `mint-consumer`: Event consumer for token minting
- **Local Environment**: Docker Compose + make targets per service

### Essential Commands

```sh
# Service-specific builds
cd services/cashback-service-api && make run
cd services/blockchain-adapter && make run
cd services/mint-consumer && make run

# Tests (per service)
make test             # Run tests in current service
```

---

## Documentation Output Rules

Do not create physical `.md` files for explanations, guides, or troubleshooting.

- Show explanations in chat (in-memory) only
- Do not create files like `SOLUTION.md`, `TROUBLESHOOTING.md`, `EXPLANATION.md`

Only create `.md` files for project documentation that belongs in the repository:
- `README.md`, `ARCHITECTURE.md`, `API.md`, `CONTRIBUTING.md`

---

## Project Architecture

Clean architecture with domain-driven design:

- **Domain Layer**: Core business entities and rules — no JSON tags, no cross-domain imports
- **Use Case Layer**: Application-specific business logic — defines `Repository` and `TransactionManager` as local interfaces
- **Repository Layer**: Data access — propagates context through all database operations
- **Handler Layer**: HTTP request/response — maps errors to status codes, uses `RegisterEndpoints`

### Directory Layout

```
/services/{service-name}/
    cmd/
        main.go or api/main.go
    internal/
        app/{domain}/
            domain/
            handler/{operation}/
            repository/
            usecase/{operation}/
            errors.go
        config/
        bootstrap/
    pkg/                # Service-specific utilities
/pkg/                   # Shared utilities across services
/db/                    # Database schema files
/proto/                 # Protocol buffer definitions
```

### Key Rules (full detail in rules files)

| Rule | File |
|------|------|
| Clean architecture, new endpoint guide, layer contracts | `.github/ai/rules/codebase/rule-1-clean-architecture.md` |
| Transaction pattern (TransactionManager interface) | `.github/ai/rules/codebase/rule-2-transaction-pattern.md` |
| Go style (grouping, no else, any, errors.Is, naming) | `.github/ai/rules/language/go/rule-3-go-style.md` |
| Testing (should pattern, table-driven, mockery, build tags) | `.github/ai/rules/language/go/rule-4-testing.md` |
| Error handling (never log+return, domain errors, wrapping) | `.github/ai/rules/language/go/rule-5-error-handling.md` |

---

## AI Skills

Use the **Orchestrator** as the single entry point for any task involving codebase analysis or changes. Do not perform code analysis, implementation, or review directly — always delegate through the appropriate skill.

| Task | Skill |
|------|-------|
| Any codebase analysis or research | `orchestrator/SKILL.md` → delegates to `researcher/SKILL.md` |
| Planning an implementation | `orchestrator/SKILL.md` → delegates to `planner/SKILL.md` |
| Implementation | `orchestrator/SKILL.md` → delegates to `coder/SKILL.md` |
| Code review | `reviewer/SKILL.md` |
| Writing and executing commits | `git-committer/SKILL.md` |
| Compress session to avoid context loss | `context-compressor/SKILL.md` (also: `/compress`) |

All skill files are under `.github/ai/skills/`.

### Plans are local-only

Plans live in `~/ai-plans/{repo-name}/{slug}/` and are **never committed**. They are stored
completely outside the repository — immune to all git operations.

```bash
# Path for this repo
~/ai-plans/cashback-platform/{slug}/
```

- Plans survive branch switches, stash, rebase, and `git clean` (they are outside the repo)
- Each developer owns their own local plans
- To resume on a different machine: run `/compress` first to generate `session-summary.md`,
  then share that file manually and use the re-attach prompt at the bottom of it

### Plan Discovery

Skills discover the active plan automatically. When no slug is explicit, they scan
`~/ai-plans/{repo-name}/` and read `progress.md` in each directory to find `IN_PROGRESS` plans.
If a ticket ID is provided, the slug is matched by prefix.

```bash
# Path for this repo
~/ai-plans/cashback-platform/{slug}/
```

- Plans survive branch switches, stash, rebase, and `git clean` (they are outside the repo)
- Each developer owns their own local plans
- To resume on a different machine: run `/compress` first to generate `session-summary.md`,
  then share that file manually and use the re-attach prompt at the bottom of it

### Critical: Codebase Analysis

Any analysis that involves reading source files, tracing call chains, identifying impacted paths, or researching existing behavior must go through the Orchestrator → Researcher flow. Do not answer analysis questions from memory or partial context — always read the relevant files first.

---


## Database Migrations

- Database schema is managed in `db/schema.sql`
- Apply migrations manually before testing changes

---

## Observability

- Always propagate `ctx` through all calls for trace ID support
- Use `internal/logger` for structured logging
- Use New Relic helpers for custom events on critical paths

---

## PR and Branching Conventions

### PR Checklist

- [ ] Tests added/updated
- [ ] Database schema changes reviewed
- [ ] All relevant modules wired
- [ ] Documentation updated if needed

### PR Description Format

- **Summary**: What and why
- **New Components**: New modules/domains added
- **Refactoring**: Code reorganization
- **New Types**: New types/interfaces introduced

### PR Title Format

Follow your project's convention. Common patterns:
- `Brief description of the change`
- `TICKET-123 | Brief description`
- `[component] Brief description`

### Branching Strategy

Adapt to your team's workflow. Common patterns:
- Feature branches from `main` or `develop`
- Bugfix branches for releases
- Hotfix branches for production issues
- Release branches for version management

---

## Commit Conventions

Full execution skill: `.github/ai/skills/git-committer/SKILL.md`

### Format

Adapt to your project's convention. Common patterns:

**With ticket tracking:**
```
<ticket-id> | <subject>

<body>
```

**Without ticket tracking:**
```
<subject>

<body>
```

### General Guidelines

- **Subject line**: Keep under 50-72 characters, imperative mood, capitalized, no period
- **Body**: Wrap at 72 characters, explain what and why (not how)
- **Separation**: Blank line between subject and body

### Chris Beams' Seven Rules (recommended)

1. Separate subject from body with a blank line
2. Limit the subject line to 50 characters
3. Capitalize the subject line
4. Do not end the subject line with a period
5. Use imperative mood — "Add feature" not "Added" / "Adding"
6. Wrap the body at 72 characters
7. Use the body to explain **what** and **why**, not how

### Example

```
Add consumer health check endpoint

Implements HTTP /health endpoint that verifies consumer connectivity
to message broker and database.
```

