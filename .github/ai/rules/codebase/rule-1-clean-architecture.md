# Rule 1 — Clean Architecture & Project Conventions

## Layer Overview

```
internal/app/{domain}/
├── domain/          → Core entities, types, constants — no external deps, no json tags
├── usecase/{op}/    → Business logic, repository interface, TransactionManager interface
├── repository/      → Data access implementation
└── handler/{op}/    → HTTP parse → use case → response
```

## Critical Layer Rules

### Domain
- No JSON tags on domain structs — domain must not depend on external layers
- No cross-domain imports — leads to cyclic dependencies
- Define typed constants/enums over magic strings

### Use Case
- Define `Repository` interface locally — do NOT import repository package
- Define `TransactionManager` as a local interface if transactions are needed
- For simple CRUD: domain types as input/output is acceptable
- For custom output: define input/output types within the use case package itself

```go
type (
    Repository interface {
        Save(ctx context.Context, entity domain.Entity) error
    }

    // Define locally if transactions are needed
    TransactionManager interface {
        WithTransaction(ctx context.Context, cb func(ctx context.Context) error) error
    }

    UseCase struct {
        repository         Repository
        transactionManager TransactionManager // optional, if needed
    }
)
```

### Repository
- Use context-aware database access pattern
- Propagate `ctx context.Context` through all database operations
- Convert database models to domain entities via `ToDomain()` methods
- Encapsulate database-specific logic (GORM, SQL, etc.)

```go
type Repository struct {
    db DatabaseInterface // Your database abstraction
}

func (r Repository) Save(ctx context.Context, entity domain.Entity) error {
    model := toModel(entity)
    return r.db.WithContext(ctx).Create(&model).Error
}
```

### Handler
- Map errors to HTTP status codes via error mapping
- Parse request → call use case → format response
- Register endpoints following project routing conventions
- Use structured logging with trace/request IDs

### Module Wiring (Dependency Injection)
- Wire concrete implementations here
- Pass database connections, transaction managers, external clients
- This is where concrete types meet interface contracts

```go
// Example wiring
repo := repository.NewRepository(db)
txManager := transaction.NewManager(db) // If using transactions
uc := usecase.NewUseCase(repo, txManager)
handler := handler.NewHandler(uc)
```

## Endpoint Registration

- Register endpoints in handler layer
- Wire dependencies in bootstrap/module layer
- Follow project-specific routing and middleware patterns
