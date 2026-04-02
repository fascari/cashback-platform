---
applyTo: "**/*.go"
---

# Package Design Philosophy

Based on Bill Kennedy's [Design Philosophy on Packaging](https://www.ardanlabs.com/blog/2017/02/design-philosophy-on-packaging.html) and [Package Oriented Design](https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html), adapted to this project's structure.

## Core Principle

> Packages must provide, not contain. A package exists to solve a specific problem domain, not to be a dumping ground of disparate concerns.

If you struggle to name a package, it is not focused on a single purpose. If a package name describes what it contains rather than what it provides, it needs restructuring.

## Project Package Layers

This project maps Bill Kennedy's layers as follows:

| Kennedy's layer | Project equivalent | Location | Purpose |
|---|---|---|---|
| Kit | Shared utilities | `pkg/` | Foundational, high-portability packages |
| internal/platform | Infrastructure | `internal/` (non-app) | Project-specific foundational support |
| internal | Business domains | `internal/app/{domain}/` | Domain logic, CRUD, services |
| cmd | Programs | `cmd/`, `consumers/cmd/`, `jobs/cmd/` | Entry points, startup, shutdown |

## Package Location Rules

### `pkg/` — Shared Utilities (Kit)

Highest portability. Could theoretically be extracted to a separate module.

- MUST NOT import from `internal/` or `cmd/`
- MUST NOT set policy about application concerns (logging framework, specific DB, config format)
- MUST NOT panic
- MUST NOT wrap errors: return only root cause values
- Configuration and runtime changes must be decoupled (accept interfaces, not concrete config)
- Each package provides a single, focused capability: `clock`, `filter`, `paginator`, `money`, `slices`

### `internal/` (non-app) — Infrastructure

Project-specific foundational support: `dbtx`, `audit`, `logger`, `tokenmanager`, `middleware`.

- MUST NOT import from `cmd/` or `internal/app/`
- MUST NOT panic
- MUST NOT wrap errors: return only root cause values
- CAN import from `pkg/`
- CAN import from other `internal/` infrastructure packages

### `internal/app/{domain}/` — Business Domains

Domain-specific business logic. Each domain is a self-contained unit.

- MUST NOT import from `cmd/`
- MUST NOT import from other `internal/app/` domains at the same level
- CAN import from `internal/` infrastructure and `pkg/`
- CAN wrap errors with context
- CAN set application-level policy (logging, config)
- Sub-packages within a domain can import each other

### `cmd/`, `consumers/cmd/`, `jobs/cmd/` — Programs

Entry points that wire everything together.

- CAN import from any internal or pkg package
- CAN panic on fatal startup errors
- CAN recover panics
- This is where most error handling (logging) happens
- This is where concrete types are wired to interfaces

## Forbidden Package Names

Never create packages named:

- `utils`, `helpers`, `common`, `misc`, `shared`, `base`, `core`, `types`, `model`

These names describe what a package contains, not what it provides. They become dumping grounds for disparate concerns and create project-wide coupling.

Instead, name packages by what they provide:

| Bad | Good | Reason |
|-----|------|--------|
| `utils/money.go` | `money/` | Provides monetary operations |
| `helpers/filter.go` | `filter/` | Provides filtering capabilities |
| `common/clock.go` | `clock/` | Provides time abstraction |
| `types/pagination.go` | `paginator/` | Provides pagination logic |

## Dependency Direction

Dependencies flow inward. Outer layers depend on inner layers, never the reverse:

```
cmd/ ──→ internal/app/ ──→ internal/ ──→ pkg/
                │                │
                └────────────────┘
```

Cross-domain imports within `internal/app/` are forbidden. If two domains need to share logic:

1. Extract the shared concept into an `internal/` infrastructure package
2. Or define interfaces at the consumer side and let the wiring layer connect them

## Preventing Single Points of Dependency

A "common types" package that many packages depend on creates project-wide coupling. Any change to it can break unrelated packages. Instead:

- Define types close to where they are used
- Use interfaces to decouple dependencies
- Accept the cost of some duplication over the cost of tight coupling
- Domain types belong in `internal/app/{domain}/domain/`, not in a shared package

## Validate Package Design

When creating or reviewing a package, verify:

1. **Location**: Does the package belong where it is? Does it follow the import rules for its layer?
2. **Dependencies**: Does it import only from allowed layers? Are cross-level imports justified?
3. **Policy**: Does a `pkg/` or infrastructure package avoid setting application policy?
4. **Naming**: Does the name describe what the package provides, not what it contains?
5. **Focus**: Does the package solve one specific problem domain?
6. **Portability**: Could the package be used in another project without modification? (for `pkg/` packages)
