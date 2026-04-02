---
applyTo: "**/*.go"
---

# Error Handling

## Core Principle

> Never log and return the same error. Choose one: log it OR return it.

## Domain Error Codes

Define error code constants in `errors.go` at the domain package root. Use `apperror.AppError` for errors that carry a code and message:

```go
// internal/app/{domain}/errors.go
package domain

const (
    ErrCodeInvalidStatusTransition = "error_invalid_status_transition"
    ErrCodeEntityNotFound          = "error_entity_not_found"
    ErrCodeBusinessValidation      = "error_business_validation"
)
```

Create errors with `apperror.New`:

```go
import "github.com/your-org/your-project/pkg/apperror"

return apperror.New(domain.ErrCodeEntityNotFound, "entity %s not found", id)
```

Check errors with `apperror.As`:

```go
if apperror.As(err, domain.ErrCodeEntityNotFound) {
    // handle not found
}
```

## Struct-Based Domain Errors

For errors that carry typed data beyond code+message, use struct errors with `Is()`:

```go
type (
    ErrNotFound     struct{ ID string }
    ErrInvalidInput struct{ Message string }
)

func (e ErrNotFound) Error() string { return fmt.Sprintf("entity not found: %s", e.ID) }
func (e ErrNotFound) Is(target error) bool {
    _, ok := target.(ErrNotFound)
    return ok
}
```

Compare with `errors.Is()`, never `==`:

```go
if errors.Is(err, ErrNotFound{}) { ... }
```

## Error Wrapping

Always wrap with context describing the operation that failed:

```go
return fmt.Errorf("finding user by id %s: %w", id, err)
```

## Handler Error Mapping

Every domain handler package has an `error_mapping.go` that maps domain error codes to HTTP status codes.

### Structure

File location: `internal/app/{domain}/handler/error_mapping.go`

```go
package handler

import (
    "net/http"

    "github.com/your-org/your-project/internal/app"
    "github.com/your-org/your-project/internal/app/{domain}/domain"
    "github.com/your-org/your-project/pkg/http/errorhandler"
)

var (
    ErrorMapping = errorhandler.ErrorMapping{
        // Domain-specific errors
        domain.ErrCodeEntityNotFound:          http.StatusNotFound,
        domain.ErrCodeBusinessValidation:      http.StatusUnprocessableEntity,
        domain.ErrCodeInvalidStatusTransition: http.StatusUnprocessableEntity,

        // Shared app errors
        app.ErrCodeBadRequest:    http.StatusBadRequest,
        app.ErrCodeRecordNotFound: http.StatusNotFound,
    }
)
```

### How Handlers Use It

Handlers call `errorhandler.Render` with the domain's `ErrorMapping`:

```go
func (h Handler) Handle(c *gin.Context, w http.ResponseWriter, r *http.Request) error {
    ctx := c.Request.Context()

    result, err := h.useCase.Execute(ctx, input)
    if err != nil {
        return errorhandler.Render(ctx, w, err, errorhandler.WithErrorMapping(handler.ErrorMapping))
    }

    return web.EncodeJSON(w, response, http.StatusOK)
}
```

### Resolution Order

`errorhandler.Render` resolves the HTTP status code in this order:

1. Built-in defaults: `RecordNotFound` → 404, `DuplicateEntry` → 409
2. Domain-specific mapping from the `ErrorMapping` map
3. Fallback: `500 Internal Server Error`

### HTTP Status Mapping Conventions

| Error category | HTTP status | When to use |
|---|---|---|
| Entity not found | `404 Not Found` | Resource does not exist |
| Validation failure | `400 Bad Request` | Invalid input, missing fields, bad format |
| Business rule violation | `422 Unprocessable Entity` | Valid syntax but violates domain rules |
| State conflict | `409 Conflict` | Duplicate, concurrent modification, invalid state transition |
| Audit/system failure | `500 Internal Server Error` | Missing request context, system errors |

### JSON Error Response

All errors produce a consistent JSON response:

```json
{
    "code": "error_entity_not_found",
    "message": "entity abc-123 not found"
}
```

Logging is handled by `errorhandler.Render`: 5xx errors log at ERROR level, 4xx at WARN.

### Adding a New Domain Error

1. Define the error code constant in `internal/app/{domain}/errors.go`
2. Use `apperror.New(code, message)` where the error is raised
3. Add the mapping entry in `internal/app/{domain}/handler/error_mapping.go`
4. The handler already passes `ErrorMapping` to `errorhandler.Render`, so new codes are picked up automatically

## What NOT To Do

```go
// Bad: log AND return
log.Error("failed to find entity", "error", err)
return err

// Good: return only (errorhandler.Render logs it)
return fmt.Errorf("finding entity: %w", err)

// Good: OR log only (fire and forget, e.g., background jobs)
log.Error("background job failed", "error", err)
```
