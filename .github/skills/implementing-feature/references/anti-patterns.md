# Anti-Patterns Reference

Recurring mistakes to check before presenting code. This file contains project-specific anti-patterns.

> **For new projects**: replace the examples below with real patterns discovered in this codebase. Update this file when new recurring mistakes are found during code review.

## Architecture

| Category | Always do | Never do |
|----------|-----------|----------|
| Domain models | Keep domain/business types free of external layer concerns (no serialization tags, no ORM annotations) | Add JSON/XML tags or ORM annotations to domain types |
| Cross-module imports | Redefine types locally at the boundary; a little copying > a little dependency | Import one domain's internal types inside another domain |
| Service dependencies | Define dependency interfaces locally in the service/use-case | Import concrete infrastructure packages directly in service code |

## Data Access

| Category | Always do | Never do |
|----------|-----------|----------|
| DB access | Use the context-aware DB accessor so transactions propagate via context | Bypass the context and use a raw DB handle directly |
| Transactions | Use the transaction manager abstraction | Use the ORM's raw transaction API directly in service code |

## Testing

| Category | Always do | Never do |
|----------|-----------|----------|
| Test data | Use factory/fixture helpers | Inline complex structs directly in test bodies |
| Assertions | Use fail-fast assertions (test stops on first failure) | Use soft assertions that continue after failure — they cascade panics |
| Mock setup | Use the project's `EXPECT()`-style builder | Use string-based `On("MethodName", ...)` setup |

## Pattern References

Add references to specific files in this codebase that demonstrate correct patterns:

- Domain type example: `{path/to/example}`
- Service/use-case example: `{path/to/example}`
- Data access example: `{path/to/example}`
- Test factory example: `{path/to/example}`

