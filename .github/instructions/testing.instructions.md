---
applyTo: "**/*_test.go,**/testdata/**/*.go,**/factory/**/*.go"
---

# Testing

## Test Naming Convention

```go
func TestUseCaseName_ShouldDescribeExpectedBehavior(t *testing.T) { ... }
// subtests: "should return error when entity not found"
```

## Table-Driven vs Individual Tests

Use table-driven tests when:

- Multiple cases share the same setup/act/assert structure and differ only in inputs and expected outputs
- You are testing a pure function or a method with a single code path that branches on input
- The test cases are independent and do not require unique setup logic per case

Use individual test functions when:

- Each case requires significantly different setup (different mocks, different dependencies, different state)
- The test exercises a complex flow where the arrange step is the main substance of the test
- A table would need fields like `setupFunc`, `mockBehavior`, and `customAssert` for most rows, making the table harder to read than separate functions
- You are testing distinct behaviors that happen to live in the same function (e.g., create vs update path in a single `Save` method)

The goal is readability. A table with 3 simple rows is clearer than 3 separate functions. A table with 15 rows where each row has a unique `setup` closure and custom assertions is harder to follow than 15 named functions.

```go
// Good: table-driven, cases differ only in input/output
func TestCalculateDiscount_ShouldApplyCorrectRate(t *testing.T) {
    tests := []struct {
        name     string
        quantity int
        want     float64
    }{
        {name: "should return zero for small orders", quantity: 5, want: 0},
        {name: "should return 10% for medium orders", quantity: 50, want: 0.10},
        {name: "should return 15% for large orders", quantity: 200, want: 0.15},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := CalculateDiscount(tt.quantity)
            require.Equal(t, tt.want, got)
        })
    }
}

// Good: individual tests, each requires distinct setup
func TestPublishPost_ShouldPublishWhenValid(t *testing.T) {
    repo := mocks.NewRepository(t)
    repo.EXPECT().FindByID(mock.Anything, "post-1").Return(draftPost(), nil)
    repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)
    notifier := mocks.NewNotifier(t)
    notifier.EXPECT().Notify(mock.Anything, mock.Anything).Return(nil)

    err := NewUseCase(repo, notifier).Execute(ctx, "post-1")
    require.NoError(t, err)
}

func TestPublishPost_ShouldRejectWhenAlreadyPublished(t *testing.T) {
    repo := mocks.NewRepository(t)
    repo.EXPECT().FindByID(mock.Anything, "post-1").Return(publishedPost(), nil)

    err := NewUseCase(repo, nil).Execute(ctx, "post-1")
    require.True(t, errors.Is(err, ErrAlreadyPublished))
}
```

## Table-Driven Tests

```go
func TestUseCase_ShouldReturnEntity(t *testing.T) {
    tests := []struct {
        name    string
        setup   func(*mocks.Repository)
        wantErr error
    }{
        {
            name: "should return entity when found",
            setup: func(m *mocks.Repository) {
                m.EXPECT().FindByID(mock.Anything, "123").Return(domain.Entity{ID: "123"}, nil)
            },
        },
        {
            name: "should return error when not found",
            setup: func(m *mocks.Repository) {
                m.EXPECT().FindByID(mock.Anything, "999").Return(domain.Entity{}, ErrNotFound)
            },
            wantErr: ErrNotFound,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := mocks.NewRepository(t)
            tt.setup(repo)
            got, err := NewUseCase(repo).Execute(context.Background(), tt.input)
            if tt.wantErr != nil {
                require.Error(t, err)
                require.True(t, errors.Is(err, tt.wantErr))
                return
            }
            require.NoError(t, err)
            require.Equal(t, tt.want, got)
        })
    }
}
```

## Mock Generation

```go
//go:generate mockery --all --case=snake --disable-version-string --with-expecter
```

- Always use `EXPECT()` builder: **never** `mock.On()`
- Use `mock.Anything` for context parameters

## Assertions

- `require` for ALL assertions: stops test on failure
- **Never** use `assert`: test continues after failure and cascades panics

```go
// Correct
require.NoError(t, err)
require.Equal(t, want, got)
require.True(t, errors.Is(err, tt.wantErr))
require.Len(t, results, 3)
require.Empty(t, results)
```

```go
// Wrong — must never appear in any test file
assert.Equal(t, want, got)   // ← test continues even when this fails
assert.NoError(t, err)       // ← subsequent lines may panic on nil dereference
```

This rule is absolute: **no exceptions, no `assert` anywhere in test files**, including handler tests, table-driven loops, and suite subtests.

## Test Data

Never define complex test data inline in test files.

### testdata/ package (per feature)

For use case and handler tests, create `testdata/` within the package. Each file is named after the domain entity it constructs — **never by role** (`inputs.go`, `expected.go`, `errors.go` are wrong):

```
pkg/{feature}/
├── feature.go
├── feature_test.go
└── testdata/
    ├── user.go       ← all User factory functions
    ├── purchase.go   ← all Purchase factory functions
    └── cashback.go   ← all Cashback factory functions
```

Each entity file:
- Contains every factory function for that entity covering all states needed by the tests
- Groups constants (IDs, mock timestamps) with the entity that primarily owns them
- Uses no artificial split between "inputs" and "expected outputs" — both live in the same file

Function names describe the **specific state** of the entity, not just the type:

```go
// testdata/cashback.go
package testdata

import cashdomain "github.com/example/internal/app/cashback/domain"

const (
    CashbackID int64 = 1
    UserID     int64 = 10
)

func ApprovedCashback() cashdomain.Cashback {
    return cashdomain.Cashback{ID: CashbackID, UserID: UserID, Status: cashdomain.StatusApproved}
}

func PendingCashback() cashdomain.Cashback {
    return cashdomain.Cashback{ID: CashbackID, UserID: UserID, Status: cashdomain.StatusPending}
}
```

```go
// testdata/user.go
package testdata

import userdomain "github.com/example/internal/app/user/domain"

func FoundUser() userdomain.User {
    return userdomain.User{ID: UserID, Email: "user@example.com"}
}

func CreatedUser() userdomain.User {
    return userdomain.User{ID: UserID, Email: "user@example.com", WalletAddress: "0xABC"}
}
```

One file per entity is the default. Split into multiple files for the same entity (e.g., `offer.go`, `offer_phc.go`) only when the file grows large or the variants are conceptually distinct.

### internal/test/factory/ (shared across domains)

For entities reused across multiple domains, use the shared factory:

```
internal/test/factory/
├── user.go
├── order.go
└── product.go
```

## Integration Tests

Tag every file:

```go
//go:build integration
```

Use `pkg/testsuite` suite + YAML fixtures:

```go
//go:build integration

type RepositorySuite struct {
    testsuite.Suite
    repo repository.Repository
}

func TestRepositorySuite(t *testing.T) { testsuite.Run(t, &RepositorySuite{}) }

func (s *RepositorySuite) SetupSuite() {
    s.Suite.SetupSuite()
    s.Suite.ConfigureFixtures("default")
    s.repo = repository.NewRepository(s.DB)
}

func (s *RepositorySuite) TestFindByID_ShouldReturnEntityWhenExists() {
    result, err := s.repo.FindByID(context.Background(), "2a8fa59d-...")
    s.Require().NoError(err)
    s.Require().Equal("USR001", result.UserID)
}
```

### Database assertions in integration tests

Never write raw DB queries inline inside a test function. Any assertion that requires querying the database to verify side effects belongs in a dedicated `assert/` sub-package under `testdata/`:

```
internal/app/{domain}/repository/testdata/
├── fixtures/
├── inputs.go
└── assert/
    └── {entity}.go   ← database assertion helpers
```

Each function in `assert/` must:
- Accept `t *testing.T` and `db *gorm.DB` as first parameters
- Call `t.Helper()` as the first statement
- Use `require` (never `assert`) for all assertions
- Have a descriptive name that reads as a sentence: `OrderCancelled`, `UserDeactivated`

```go
//go:build integration

package assert

import (
    "testing"

    "github.com/stretchr/testify/require"
    "gorm.io/gorm"

    "github.com/your-org/your-project/internal/app/order/domain"
)

func OrderCancelled(t *testing.T, db *gorm.DB, orderID string) {
    t.Helper()
    var status string
    err := db.Raw("SELECT status FROM orders WHERE id = ?", orderID).Scan(&status).Error
    require.NoError(t, err, "should query order without error")
    require.Equal(t, string(domain.StatusCancelled), status, "order %s should be CANCELLED", orderID)
}
```

Call helpers directly from the test body or from a closure in a table field — the choice depends on whether the case fits naturally in the table:

```go
// Individual test — when the DB assertion makes this case distinct enough to stand alone
func (cs *Suite) TestCancelOrder_ShouldCancelConflictingActiveOrder() {
    result, err := cs.repository.CancelOrder(cs.ctx, testdata.OrderID)

    cs.Require().NoError(err)
    cs.Require().Equal(testdata.CancelledOrder(), result)
    orderassert.OrderCancelled(cs.T(), cs.DB, testdata.OrderID)
}
```

```go
// Table field — when the case fits naturally alongside other cases in the table
{
    name: "Should cancel active order",
    ...
    assert: func() {
        orderassert.OrderCancelled(cs.T(), cs.DB, testdata.OrderID)
    },
},
```

The non-negotiable rule is: **never inline raw DB queries in the test**. Always delegate to a helper in `assert/`. Whether you call that helper from an individual test or from a table closure is a readability judgment.



```yaml
# testdata/fixtures/default/users.yaml
- id: '2a8fa59d-cab7-47d8-ad07-c45b4d1d1279'
  user_id: USR001
  status: active
  email: user@example.com
```

## Comments in Tests

Test code is self-describing. The function name, subtest strings, and variable names are the documentation.

**Default: no comments.** This applies equally to `*_test.go` files and `testdata/` packages.

Forbidden:
- Function-level doc comments (`// TestFoo verifies that...`)
- `// Arrange / Act / Assert` section markers
- Inline comments restating assertions (`require.Empty(t, result) // should be empty`)
- Var-level doc comments on exported testdata variables

```go
// Bad
// TestFoo_ShouldReturnError verifies that an error is returned when not found.
func TestFoo_ShouldReturnError(t *testing.T) { ... }

// Good
func TestFoo_ShouldReturnErrorWhenNotFound(t *testing.T) { ... }
```

**Exception: use a comment only when ALL three conditions are true:**

1. The behavior is non-obvious and cannot be expressed by renaming
2. The comment explains WHY the test is structured that way
3. Removing the comment would force the reader to trace business logic

```go
// Good: explains non-obvious domain constraint
func TestActivateSubscription_ShouldSkipWhenTrialActive(t *testing.T) {
    // Trial subscriptions bypass the activation flow: they are already active
    // by definition and skipping prevents a duplicate activation error.
    ...
}
```

## What to Test

- Happy path
- Error cases (each error type)
- Edge cases (empty, nil, zero values)
- Transaction rollback scenarios (for write operations)
