---
applyTo: "internal/app/**/*.go"
---

# Clean Architecture & Domain-Driven Design

## Domains

A domain is a bounded context that owns a specific business capability. Each domain is a self-contained unit with its own entities, rules, errors, and data access.

### What defines a domain

- It represents a real business concept: `order`, `inventory`, `customer`, `notification`
- It owns its data: the domain decides its table structure, validation rules, and state transitions
- It has a clear API boundary: other parts of the system interact with it through its use cases, never by reaching into its internals
- It is independently testable: a domain can be unit-tested without standing up other domains

### Domain boundaries are sacred

- **No cross-domain imports at the same level.** `order` must not import from `inventory`. If two domains need to coordinate, use events, a shared infrastructure package, or orchestrate from the wiring layer
- **Shared concepts go to infrastructure.** If two domains need the same type (e.g., money, audit, pagination), it belongs in `internal/` or `pkg/`, not duplicated or cross-imported
- **Domain logic never leaks into handlers or repositories.** Validation, state transitions, and business rules live in the domain or use case layer. Repositories only persist and retrieve; handlers only parse and respond

### When to create a new domain

Create a new domain when:
- The concept has its own lifecycle (created, updated, deleted independently)
- It has its own business rules and validation
- It maps to a distinct set of database tables
- A different team member could own it without needing to understand other domains

Do NOT create a new domain when:
- The concept is just a supporting detail of an existing domain (use a sub-package instead)
- It would only have a single entity with no behavior

## Layer Overview

```
internal/app/{domain}/
├── domain/          → Core entities, types, constants, business rules
├── usecase/{op}/    → Application logic, local interfaces for Repository and TransactionManager
├── repository/      → GORM + BaseRepository, reader.go / writer.go split
└── handler/{op}/    → HTTP parse → use case → response
```

## Data Flow Between Layers

Domain types are the lingua franca of the application. They flow freely between layers: handlers receive them from use cases, use cases receive them from repositories. However, each layer has its own representation concerns:

- **Domain types** (`domain/`): pure business representation. No `json`, `gorm`, or any framework tags. These are the types that all layers share
- **Models** (`repository/`): GORM-annotated structs with `gorm:` tags. Models are internal to the repository and MUST NEVER leak outside it. Repositories accept domain types, convert to models internally, persist, and convert back to domain types via `ToDomain()`
- **DTOs / Payloads** (`handler/`): request/response structs with `json:` tags. Handlers convert between DTOs and domain types. DTOs are internal to the handler

```
Handler (DTO ↔ domain) → UseCase (domain) → Repository (domain ↔ model)
```

### What must never cross a boundary

| Type | Belongs to | Must never appear in |
|------|-----------|---------------------|
| Model (gorm tags) | `repository/` | use case, handler, domain |
| DTO (json tags) | `handler/{op}/` | use case, repository, domain |
| Domain type | `domain/` | everywhere (this IS the shared type) |

```go
// WRONG: domain struct with json tags
type User struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

// CORRECT: domain struct, no framework tags
type User struct {
    ID   string
    Name string
}
```

```go
// WRONG: repository returning a model
func (r Repository) FindByID(ctx context.Context, id string) (model, error)

// CORRECT: repository converts model to domain internally
func (r Repository) FindByID(ctx context.Context, id string) (domain.User, error) {
    var m userModel
    if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
        return domain.User{}, err
    }
    return m.ToDomain(), nil
}
```

## Domain Layer

- No `json`, `gorm`, or any framework tags on domain structs
- No cross-domain imports: leads to cyclic dependencies
- Define typed constants/enums over magic strings
- Domain errors live in `errors.go` at the domain package root
- Business rules and validations belong here, not in use cases or handlers

```go
type (
    Status string
)

const (
    StatusDraft  Status = "DRAFT"
    StatusActive Status = "ACTIVE"
)
```

## Use Case Layer

- Define `Repository` interface locally: do NOT import the repository package
- Define `TransactionManager` as a local interface: NEVER import `internal/dbtx`
- For simple CRUD: domain types as input/output is acceptable
- For custom output: define input/output types within the use case package itself
- Use cases orchestrate domain logic and repository calls; they do not contain business rules themselves

```go
type (
    Repository interface {
        Save(ctx context.Context, entity domain.Entity) error
    }

    TransactionManager interface {
        WithTransaction(ctx context.Context, cb func(ctx context.Context) error) error
    }

    UseCase struct {
        repository         Repository
        transactionManager TransactionManager
    }
)
```

## Repository Layer

- NEVER use `db.Transaction()` or `db.WithTransaction()` from GORM directly
- Use `ToDomain()` conversion methods on models
- Always pass `ctx` to GORM via `.WithContext(ctx)` for trace propagation

### Simple repository (no distributed transactions)

When the use case performs single logical writes, the repository holds a plain `*gorm.DB` field. There is no transaction to join, so `dbtx.BaseRepository` is unnecessary overhead.

```go
type Repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
    return Repository{db: db}
}

func (r Repository) Save(ctx context.Context, entity domain.Entity) error {
    return r.db.WithContext(ctx).Create(&model).Error
}
```

### Transaction-aware repository (distributed transactions)

When the use case needs `TransactionManager` to wrap multiple repository calls atomically, the repository MUST embed `dbtx.BaseRepository` and use `r.DB(ctx)` so it joins the transaction propagated via context.

```go
type Repository struct {
    dbtx.BaseRepository
}

func NewRepository(db *gorm.DB) Repository {
    return Repository{BaseRepository: dbtx.NewBaseRepository(db)}
}

func (r Repository) Save(ctx context.Context, entity domain.Entity) error {
    return r.DB(ctx).Create(&model).Error
}
```

### Decision table

| Use case has `TransactionManager`? | Repository embeds | DB access | Why |
|---|---|---|---|
| No | plain `*gorm.DB` | `r.db.WithContext(ctx)` | No transaction to join; simpler and correct |
| Yes | `dbtx.BaseRepository` | `r.DB(ctx)` | Must join the transaction propagated via context |

> `r.DB(ctx)` checks the context for a transaction injected by `TransactionManager.WithTransaction`. Without that, it falls back to a plain connection — functionally equivalent to `r.db` but with unnecessary indirection.

## Handler Layer

- Map errors to HTTP status codes via `error_mapping.go`
- Parse request → call use case → format response
- Register endpoints via `RegisterEndpoints(r *web.RouteGroup, h Handler)`
- Use `logger.WithTraceID(ctx)` for structured logging

### Handler must use the concrete use case type — never define a UseCase interface

Handlers hold the concrete `usecase.UseCase` struct directly. **Never** define a `UseCase` interface inside the handler package.

```go
// CORRECT: import and store the concrete type
import usecase "github.com/your-org/your-project/internal/app/{domain}/usecase/{operation}"

type Handler struct {
    useCase usecase.UseCase
}

func NewHandler(useCase usecase.UseCase) Handler {
    return Handler{useCase: useCase}
}
```

```go
// WRONG: interface defined in the handler only to enable mocking
type UseCase interface {
    Execute(ctx context.Context, ...) (domain.Entity, error)
}

type Handler struct {
    useCase UseCase
}
```

A `UseCase` interface in the handler package has no architectural value — the use case is never substituted with another implementation at runtime. Its only purpose is to allow test doubles, which is the wrong layer to mock. Use case interfaces belong in the use case layer (`//go:generate mockery`), and handler tests instantiate a **real use case** with mocked repositories.

#### Handler test pattern

```go
// Correct: handler test uses real use case + mocked repositories
userRepo := ucmocks.NewUserRepository(t)
userRepo.EXPECT().FindByID(mock.Anything, "user-1").Return(user, nil)

h := getuser.NewHandler(getusersuc.NewUseCase(userRepo))
```

```go
// Wrong: handler test mocks the use case itself
mockUseCase := mocks.NewUseCase(t)  // ← mock generated from a handler-local interface
mockUseCase.EXPECT().Execute(...).Return(...)
h := getuser.NewHandler(mockUseCase)
```

## Module Wiring (DI Layer Only)

When the use case needs transactions, pass concrete `dbtx.TransactionManager` here. This is the ONLY place where `dbtx` is imported outside the repository:

```go
// With distributed transactions
repo := repository.NewRepository(db)
txManager := dbtx.NewTransactionManager(db)
uc := usecase.NewUseCase(repo, txManager)
```

When the use case does not need transactions, the wiring is simpler:

```go
// Without transactions
repo := repository.NewRepository(db)
uc := usecase.NewUseCase(repo)
```

## Uber/FX Application Structure

All services MUST follow the bootstrap + modules pattern for dependency injection wiring.

### Structure

```
cmd/{entrypoint}/
├── main.go              # Only references bootstrap.* and modules.*
└── modules/
    ├── {domain}.go      # One file per business domain
    └── types.go         # Shared FX parameter types (RouterParams etc.)

internal/bootstrap/
├── config.go            # fx.Module("config", ...)
├── database.go          # fx.Module("database", ...)
├── logger.go            # Logger() fx.Option function
├── router.go            # fx.Module("router", ...)
├── server.go            # fx.Module("server", ...)
└── {infra}.go           # One file per infra concern (nats, redis, grpc, outbox...)
```

### main.go rules

- MUST NOT contain any `fx.Provide(...)` or `fx.Invoke(...)` calls
- MUST ONLY reference `bootstrap.*` and `modules.*`
- Infrastructure logger init happens inside `bootstrap.Logger()` via `init()`

```go
// Correct
func main() {
    fx.New(
        bootstrap.Logger(),
        bootstrap.Config,
        bootstrap.Database,
        bootstrap.NATS,
        bootstrap.Server,
        modules.User,
        modules.Purchase,
    ).Run()
}

// Wrong — inline fx.Provide in main.go
func main() {
    fx.New(
        fx.Provide(nats.NewNATSClient),   // ← must be in bootstrap.NATS
        fx.Provide(config.LoadOutbox),    // ← must be in bootstrap.Config
        modules.User,
    ).Run()
}
```

### bootstrap package rules

- Each bootstrap file exports one `fx.Module("name", ...)` variable (or a `func() fx.Option` for the logger)
- Bootstrap files wire infrastructure only: config, database, cache, message broker, HTTP/gRPC server, external API clients
- Business domain logic NEVER goes in bootstrap files

```go
// bootstrap/nats.go
var NATS = fx.Module("nats",
    fx.Provide(nats.NewNATSClient),
)
```

### modules package rules

Each domain module file exports one `fx.Options(...)` variable named after the domain:

```go
// cmd/modules/user.go
var User = fx.Options(
    userFactories,    // fx.Provide(repo, use cases, handlers)
    userDependencies, // fx.Provide(interface adapter functions)
    userInvokes,      // fx.Invoke(register endpoints)
)
```

The three-variable pattern:
- `{domain}Factories` — `fx.Provide(repo.New, usecase.New, handler.NewHandler)`
- `{domain}Dependencies` — `fx.Provide(func(repo X) usecase.SomeInterface { return repo })`
- `{domain}Invokes` — `fx.Invoke(func(params RouterParams, h handler.Handler) { handler.RegisterEndpoint(...) })`

### RouterParams

Use `fx.In` tagged structs to inject named router values:

```go
type RouterParams struct {
    fx.In
    Router    *chi.Mux   `name:"main"`
    APIRouter chi.Router `name:"api"`
}
```
