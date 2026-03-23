# Rule 3 — Go Style

Follow the [Google Go Style Guide](https://google.github.io/styleguide/go/) 100%.

## Declaration Grouping

Always group related declarations:

```go
type (
    Status string
    Entity struct { ... }
)
const (
    StatusDraft   Status = "DRAFT"
    StatusActive  Status = "ACTIVE"
)
var (
    ErrNotFound = errors.New("not found")
)
```

## Control Flow — No `else`

```go
// ✅ Early return
if err != nil {
    return err
}
doHappyPath()

// ❌ Never
if err != nil {
    return err
} else {
    doHappyPath()
}
```

## Naming

- Lowercase package names — single word preferred (`handler`, `usecase`, `repository`)
- No `Get`/`Set` prefixes for methods
- Use `any` instead of `interface{}`
- Acronyms all-caps: `userID`, `httpClient`, `apiURL`
- No `utils`, `helpers`, `common`, `misc` packages — use meaningful names: `filter`, `paginator`, `money`, `clock`

## Interfaces

- Define interfaces at the point of use (use case), not at the implementation
- Small, focused interfaces — prefer single-method over large contracts
- Accept interfaces, return structs:

```go
// ✅
func NewUseCase(repo Repository) UseCase { return UseCase{repo: repo} }

// ❌ Don't return interfaces unnecessarily
func NewUseCase() Repository { return &impl{} }
```

- Use `//go:generate mockery --all --case=snake --disable-version-string --with-expecter`

## Error Comparison

```go
// ✅
if errors.Is(err, ErrNotFound) { ... }

// ❌
if err == ErrNotFound { ... }
```

## Immutability & Value Receivers

Prefer value receivers. Pointers ONLY when:
1. Struct contains a `sync.Mutex` or must be mutated by design
2. Struct > 64 bytes AND copied frequently
3. `nil`/absence must be represented semantically

```go
// ✅ Value receiver
func (h Handler) Handle(c *gin.Context) error { ... }

// ✅ Return new value instead of mutating
func (c Config) WithTimeout(t int) Config { c.Timeout = t; return c }

// ❌ Avoid pointer just to "avoid allocation"
func ProcessData(data *Data) error { ... }
```

## Pure Functions

Same input → same output, no side effects. Prefer over impure functions.

```go
// ✅ Pure
func calculateTotal(items []Item) float64 {
    var total float64
    for _, item := range items { total += item.Price }
    return total
}

// ❌ Impure — mutates external state
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
- Comments must be simple, objective, and in English

```go
// ❌ Obvious
// FindByID finds entity by ID
func (r Repository) FindByID(ctx context.Context, id string) (Entity, error)

// ✅ Explains non-obvious business rule
// Apply institutional discount only for orders > 100 units (legacy rule from 2019 contract)
if product.Quantity > 100 { basePrice *= 0.85 }
```

## Concurrency

```go
// Channels for coordination, mutexes for shared state
func (c Consumer) run(ctx context.Context) {
    for {
        select {
        case <-ctx.Done(): return
        case msg := <-c.messages: c.handle(ctx, msg)
        }
    }
}
```

## Linting

Must pass: `golangci-lint` with `revive`, `staticcheck`, `gofumpt`, `errcheck`, `ineffassign`, `gocyclo`
