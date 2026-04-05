---
name: testing-implementation
description: Use when implementing-feature delegates testing, when validating test coverage for a feature, or when any prompt mentions tests, test files, assertions, mocks, or test adjustments
---

# Testing Implementation

Writes and executes tests, validates coverage, and ensures the implementation meets
the success criteria defined in the plan. Works after implementing-feature completes each phase.

## When to use

- implementing-feature delegates testing after a phase is complete
- User asks to write tests for a specific package or feature
- User asks to validate test coverage
- **Any prompt that mentions tests, test files, test patterns, assertions, mocks, or test adjustments** — always invoke this skill first

## Steps

1. Read `.github/plans/{slug}/implementation-plan.md` to understand success criteria for the current phase.
2. Read `.github/instructions/` — apply all project-specific testing conventions (framework, patterns, mock strategy, test data).
3. Analyze existing test files for the affected packages. Identify patterns, mock setup, and factory/fixture functions.
4. Write unit tests: follow project conventions (see `.github/instructions/`). Cover happy path + each error case + edge cases. Use test data helpers/factories — never inline complex structs.
5. Write integration tests where applicable (data layer, external integrations): follow project conventions for test tagging, suites, and YAML/JSON fixtures.
6. Run **only the affected tests** — never the full suite. Check project scripts for the right commands:
   ```bash
   # Unit tests: target the specific module/package changed
   # (use project-specific test runner: go test, npm test, pytest, etc.)

   # Integration tests: use the project-specific integration test target
   # (e.g. make integration-pkg, pytest --integration, etc.)

   # Lint: scoped to the changed module
   # (use project-specific linter: golangci-lint, eslint, pylint, etc.)
   ```
   > **Never run the full suite** when only writing or modifying tests for a specific module. Target only the affected paths.
7. Update `.github/plans/{slug}/progress.md` with test results.

## Output

Update `.github/plans/{slug}/progress.md` with:

```markdown
## Test Results — Phase {N}
- Unit tests: PASS ({N} tests)
- Integration tests: PASS / SKIPPED (no local env) / FAIL
- Coverage: ~{N}%
- Lint: PASS / {issues}
```

## Test naming

Follow the project's naming convention (see `.github/instructions/`). General pattern:

```
TestSubjectName_ShouldDescribeExpectedBehavior
{ name: "should return error when id is empty" }
{ name: "should rollback transaction on save failure" }
```

## Mock pattern

Use the project's mocking strategy (see `.github/instructions/`). Example for testify/mockery (Go):

```go
// Always EXPECT() builder
repo.EXPECT().FindByID(mock.Anything, "id-1").Return(entity, nil)

// Never mock.On()
repo.On("FindByID", mock.Anything, "id-1").Return(entity, nil)
```

## Quality checklist

- [ ] Test assertions use the project's preferred assertion style (fail-fast preferred)
- [ ] Mocks use the project's preferred mock strategy
- [ ] `//go:generate` or equivalent present on all mocked interfaces (if applicable)
- [ ] Test names follow project convention (`TestFoo_ShouldDoX` / `"should do x"`)
- [ ] Test data via factory/fixture helpers — never inline complex structs
- [ ] `testdata/` files named after domain entities (`user.go`, `purchase.go`), not by role (`inputs.go`, `expected.go`)
- [ ] Factory function names describe entity state: `CreatedUser()`, `ApprovedCashback()`, `FoundPurchase()` — not just the type name
- [ ] Integration tests tagged appropriately (e.g. `//go:build integration`, pytest marker, etc.)
- [ ] All fixtures in `testdata/` or equivalent project fixture directory
