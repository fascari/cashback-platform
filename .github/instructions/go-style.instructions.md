---
applyTo: "**/*.go"
---

# Go Style

Follow the [Google Go Style Guide](https://google.github.io/styleguide/go/) as the baseline. The rules below are project-specific extensions and emphasis points.

## Naming

### Packages

- Lowercase, single word preferred: `handler`, `usecase`, `repository`, `domain`
- Multi-word packages stay unbroken in lowercase: `testsuite`, `featureflag`
- Never: `utils`, `helpers`, `common`, `misc`, `model`, `testhelper`
- Use meaningful names: `filter`, `paginator`, `money`, `clock`

### Exported Symbols

Avoid repeating the package name in the symbol:

```go
// Bad
widget.NewWidget()
db.LoadFromDatabase()

// Good
widget.New()
db.Load()
```

### Receivers

Short (1-2 letters), abbreviation of the type, consistent across all methods:

```go
func (r Repository) FindByID(ctx context.Context, id string) (Entity, error)
func (h Handler) Handle(c *gin.Context) error
func (u UseCase) Execute(ctx context.Context, input Input) (Output, error)
```

### Variables

Length proportional to scope. Omit type-like qualifiers unless disambiguation is needed:

```go
// Good
users, err := repo.FindAll(ctx)
count := len(users)

// Bad
userSlice, fetchErr := repo.FindAll(ctx)
userCount := len(userSlice)
```

### Constants & Initialisms

MixedCaps only. Never `ALL_CAPS` or `k`-prefix:

```go
const MaxRetries = 3     // Good
const MAX_RETRIES = 3    // Bad
const kMaxRetries = 3    // Bad
```

Initialisms keep consistent case: `ID`, `URL`, `HTTP`, `API`, `DB`, `SQS`, `DLQ`, `ARN`.

| Exported | Unexported |
|----------|------------|
| `UserID` | `userID` |
| `HTTPURL` | `httpURL` |
| `SQSAPI` | `sqsAPI` |
| `DLQ` | `dlq` |

### Methods

- No `Get`/`Set` prefixes: `Name()` not `GetName()`
- Use `Compute` or `Fetch` when the call is expensive or remote
- No `Get`/`Set` prefixes for methods

## Declaration Grouping

Always group related declarations:

```go
type (
    Status string
    Entity struct { ... }
)

const (
    StatusDraft  Status = "DRAFT"
    StatusActive Status = "ACTIVE"
)

var (
    ErrNotFound = errors.New("not found")
)
```

## Control Flow

No `else`. Early returns only:

```go
// Good
if err != nil {
    return err
}
doHappyPath()

// Bad
if err != nil {
    return err
} else {
    doHappyPath()
}
```

## Interfaces

- Define at the point of use (use case), not at the implementation
- Small, focused: prefer single-method over large contracts
- Accept interfaces, return structs:

```go
// Good
func NewUseCase(repo Repository) UseCase { return UseCase{repo: repo} }

// Bad
func NewUseCase() Repository { return &impl{} }
```

## Error Comparison

```go
// Good
if errors.Is(err, ErrNotFound) { ... }

// Bad
if err == ErrNotFound { ... }
```

Use `any` instead of `interface{}`.

## Immutability & Value Receivers

Prefer value receivers. Pointers ONLY when:

1. Struct contains a `sync.Mutex` or must be mutated by design
2. Struct > 64 bytes AND copied frequently
3. `nil`/absence must be represented semantically

```go
// Good: value receiver
func (h Handler) Handle(c *gin.Context) error { ... }

// Good: return new value instead of mutating
func (c Config) WithTimeout(t int) Config { c.Timeout = t; return c }
```

## Pure Functions

Same input produces same output, no side effects. Prefer over impure functions:

```go
// Good: pure
func calculateTotal(items []Item) float64 {
    var total float64
    for _, item := range items {
        total += item.Price
    }
    return total
}

// Bad: impure, mutates external state
var globalTotal float64
func addToTotal(amount float64) { globalTotal += amount }
```

## Function Size & Abstraction

- Keep functions focused on one responsibility
- Inline logic when clear and not reused
- Extract only when logic is complex OR reused across callers
- Avoid tiny functions (< 5 lines) that create unnecessary indirection

## Comments

- Only comment when explaining WHY, not WHAT
- Never add obvious/redundant comments
- All exported types and functions must have doc comments (per Google Go Style)
- Comments must be simple, objective, and in English

```go
// Bad: obvious
// FindByID finds entity by ID
func (r Repository) FindByID(ctx context.Context, id string) (Entity, error)

// Good: explains non-obvious business rule
// Apply institutional discount only for orders > 100 units (legacy rule from 2019 contract)
if product.Quantity > 100 { basePrice *= 0.85 }
```

## Doc Comments

Required for all exported names. Start with the name of the symbol:

```go
// Repository provides data access for user entities.
type Repository struct { ... }

// Execute processes the user registration for the given input.
func (u UseCase) Execute(ctx context.Context, input Input) (Output, error)
```

## Import Organization

Group imports in three blocks separated by blank lines:

```go
import (
    "context"
    "fmt"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "github.com/your-org/your-project/internal/app/user/domain"
)
```

1. Standard library
2. Third-party packages
3. Internal packages

## Concurrency

Channels for coordination, mutexes for shared state:

```go
func (c Consumer) run(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        case msg := <-c.messages:
            c.handle(ctx, msg)
        }
    }
}
```

## Linting

Must pass: `golangci-lint` with `revive`, `staticcheck`, `gofumpt`, `errcheck`, `ineffassign`, `gocyclo`.
