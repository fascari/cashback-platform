# Rule 2 — Transaction Pattern

## Pattern Overview

Clean architecture requires that business logic (use cases) remain independent of infrastructure concerns. Transactions are an infrastructure concern, so:

1. **Use cases** define what they need via interfaces
2. **Repositories** provide context-aware database access
3. **Concrete implementations** are injected at the module wiring layer

---

## The Pattern

### 1. Repository Layer — Context-Aware Database Access

Repositories must propagate context through all database operations to support transaction management:

```go
package repository

import (
    "context"
    "your-project/internal/domain"
)

type Repository struct {
    db DatabaseInterface // Your database abstraction
}

func NewRepository(db DatabaseInterface) Repository {
    return Repository{db: db}
}

// ALWAYS propagate ctx — enables transaction awareness
func (r Repository) Save(ctx context.Context, entity domain.Entity) error {
    return r.db.WithContext(ctx).Create(&entity).Error
}

func (r Repository) FindByID(ctx context.Context, id string) (domain.Entity, error) {
    var model Model
    err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
    return model.ToDomain(), err
}
```

### 2. Use Case Layer — Define Transaction Interface Locally

**NEVER** import infrastructure packages (like transaction managers) directly. Define what you need as an interface:

```go
package usecase

import "context"

type (
    // Define repository contract locally
    Repository interface {
        Save(ctx context.Context, entity domain.Entity) error
        FindByID(ctx context.Context, id string) (domain.Entity, error)
    }

    // Define transaction manager contract locally
    TransactionManager interface {
        WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
    }

    UseCase struct {
        repository         Repository
        transactionManager TransactionManager
    }
)

func (u UseCase) Execute(ctx context.Context, input Input) error {
    // Use transaction for multi-step operations
    return u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
        entity1, err := u.repository.FindByID(txCtx, input.ID)
        if err != nil {
            return err // auto rollback on error
        }
        
        entity1.Update(input.Data)
        
        if err := u.repository.Save(txCtx, entity1); err != nil {
            return err // auto rollback on error
        }
        
        return nil // auto commit on success
    })
}
```

### 3. Module Wiring — Inject Concrete Implementations

This is the **ONLY** place where concrete infrastructure types are created and injected:

```go
package modules

import (
    "your-project/internal/transaction"
    "your-project/internal/repository"
    "your-project/internal/usecase"
)

func WireMyDomain(db *sql.DB) *handler.Handler {
    // Create concrete implementations
    repo := repository.NewRepository(db)
    txManager := transaction.NewManager(db) // Concrete transaction manager
    
    // Inject into use case (as interfaces)
    uc := usecase.NewUseCase(repo, txManager)
    
    // Wire handler
    return handler.NewHandler(uc)
}
```

---

## Key Principles

### ✅ DO
- Propagate `context.Context` through all database operations
- Define infrastructure needs as interfaces in use cases
- Inject concrete implementations at wiring layer
- Return errors from transaction callbacks for auto-rollback
- Return `nil` from transaction callbacks for auto-commit

### ❌ DON'T
- Import infrastructure packages directly in use cases
- Store database connections directly in use cases
- Call `db.Begin()`, `db.Commit()`, `db.Rollback()` in repositories
- Nest transactions manually — let the transaction manager handle it

---

## Transaction Manager Implementation Example

Your project should provide a transaction manager that implements the pattern. Example interface:

```go
package transaction

import "context"

type Manager interface {
    // WithTransaction executes fn within a database transaction.
    // If fn returns an error, the transaction is rolled back.
    // If fn returns nil, the transaction is committed.
    // The context passed to fn contains the transaction state.
    WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
```

The concrete implementation depends on your database layer (GORM, sqlx, pgx, etc.) but should:
- Store transaction state in context
- Propagate that context to repositories
- Auto-commit on success
- Auto-rollback on error or panic

---

## Enforcement Checklist

- [ ] Repositories receive and propagate `ctx context.Context`
- [ ] Use cases define `TransactionManager` as a local interface
- [ ] Use cases do NOT import infrastructure packages
- [ ] Transaction manager concrete type only appears in wiring layer
- [ ] All database operations use context-aware methods
- [ ] Transaction callbacks return errors for rollback, `nil` for commit
