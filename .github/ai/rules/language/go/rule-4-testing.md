# Rule 4 — Testing

## Test Naming Convention

```go
func TestUseCaseName_ShouldDescribeExpectedBehavior(t *testing.T) { ... }
// subtests: "should return error when entity not found"
```

## Unit Tests — Table-Driven with testify + mockery

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

- Always use `EXPECT()` builder — **never** `mock.On()`
- Use `mock.Anything` for context parameters

## Assertions

- `require` for ALL assertions — stops test on failure
- **Never** use `assert` — test continues after failure and cascades panics

```go
require.NoError(t, err)
require.Equal(t, want, got)
```

## Test Data — ALWAYS use `testdata/` package

Never define test data inline in test files. Create factory functions:

```
pkg/{feature}/
├── feature.go
├── feature_test.go
└── testdata/
    ├── inputs.go     ← factory funcs for inputs
    ├── expected.go   ← expected outputs
    └── errors.go     ← shared test errors
```

```go
// testdata/inputs.go
package testdata

func NewCashback() domain.Cashback {
    return domain.Cashback{ID: "cashback-123", Status: domain.StatusActive}
}
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
    s.Suite.ConfigureFixtures("default") // loads testdata/fixtures/default/*.yaml
    s.repo = repository.NewRepository(s.DB)
}

func (s *RepositorySuite) TestFindByID_ShouldReturnEntityWhenExists() {
    result, err := s.repo.FindByID(context.Background(), "2a8fa59d-...")
    s.Require().NoError(err)
    s.Require().Equal("CB001", result.CashbackID)
}
```

YAML fixture format:
```yaml
# testdata/fixtures/default/cashback_transactions.yaml
- id: '2a8fa59d-cab7-47d8-ad07-c45b4d1d1279'
  cashback_id: CB001
  status: active
  amount: 80.00
```

## What to Test

- Happy path
- Error cases (each error type)
- Edge cases (empty, nil, zero values)
- Transaction rollback scenarios (for write operations)
