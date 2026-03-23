# Rule 5 — Error Handling

## Core Principle

> Never log and return the same error. Choose one: log it OR return it.
> — [100 Go Mistakes #52](https://100go.co/#handling-an-error-twice-52)

## Domain Error Types

```go
// internal/app/{domain}/errors.go
package domain

import "fmt"

type (
    ErrNotFound struct{ ID string }
    ErrInvalidInput struct{ Message string }
    ErrConflict struct{ ID string }
)

func (e ErrNotFound) Error() string    { return fmt.Sprintf("entity not found: %s", e.ID) }
func (e ErrInvalidInput) Error() string { return fmt.Sprintf("invalid input: %s", e.Message) }
func (e ErrConflict) Error() string    { return fmt.Sprintf("entity already exists: %s", e.ID) }
```

## Error Comparison

```go
// ✅
if errors.Is(err, ErrNotFound{}) { ... }

// For struct errors implement Is():
func (e ErrNotFound) Is(target error) bool {
    _, ok := target.(ErrNotFound)
    return ok
}
```

## Handler Error Mapping

Map domain errors → HTTP status codes in `error_mapping.go`:

```go
func toHTTPStatus(err error) int {
    var notFound ErrNotFound
    if errors.As(err, &notFound) {
        return http.StatusNotFound
    }
    var invalid ErrInvalidInput
    if errors.As(err, &invalid) {
        return http.StatusBadRequest
    }
    return http.StatusInternalServerError
}
```

## Error Wrapping

```go
// ✅ Wrap with context
return fmt.Errorf("finding cashback by id %s: %w", id, err)

// ✅ Check wrapped errors
if errors.Is(err, ErrNotFound{}) { ... }
```

## What NOT To Do

```go
// ❌ Log AND return
log.Error("failed to find entity", "error", err)
return err

// ✅ Return only (let handler log)
return fmt.Errorf("finding entity: %w", err)

// ✅ OR log only (fire and forget)
log.Error("background job failed", "error", err)
```

