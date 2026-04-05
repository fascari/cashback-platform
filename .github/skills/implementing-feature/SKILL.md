---
name: implementing-feature
description: Use when implementing a plan, feature, or bug fix phase by phase in the codebase
---

# Implementing Feature

Implements plans phase by phase, validates each step with linter and tests,
and only commits after explicit user approval.

---

## ⛔ HARD RULE — NEVER COMMIT

This skill implements code only. It NEVER executes `git commit`, `git add`, or `git push` — not even after all phases are complete. Commits require a separate, explicit user authorization via the `committing-changes` skill.

**"Go ahead and implement"** is NOT authorization to commit.
**Completing all implementation phases** is NOT authorization to commit.

Wait for the user to explicitly say: "commit", "go ahead and commit", or equivalent.

---

## When to use

- User asks to implement a plan or feature
- orchestrating-tasks delegates implementation
- User says "implement", "code this", "start coding"
- User types /implement_plan
- User says "continue", "continue implementation", "continue the plan", "continue the code"
- User says "next phase", "next step", "proceed", "go ahead", "start"
- User says "inicia", "continua", "continua a implementação", "continua o plano", "próxima fase"
- User says "faz", "implementa", "codifica", "escreve o código", "começa a implementar"
- User says "revise", "fix", "correct", "refactor" targeting specific implementation files
- User references a phase number or task from a plan (e.g. "do phase 2", "implement task 3")
- User pastes code and asks to apply a rule or pattern to it
- Any instruction that results in writing or modifying production code

## Steps

### Step 0 — Setup

> **Critical**: compilation error checks only catch syntax/type errors. They do NOT replace the linter.
> The project linter (configured in `.github/instructions/` or project config) is the only tool that enforces all
> project-level style rules. Skipping it is not allowed — not even for "trivial" changes.
> Both compilation checks (after every edit) and linting (after every phase) must pass.

1. Ensure the plans symlink exists — run setup from `.github/skills/references/plans-setup.md`

## Implementation Protocol

### Per Phase
1. Read the full plan: `.github/plans/{slug}/implementation-plan.md`
2. Read ALL files referenced in the current phase before editing
3. **Compatibility gate** — before modifying any existing signature, type, field, or contract:
   - Does this change alter an existing HTTP request/response shape visible to callers?
   - Does this change an exported Go interface that other packages implement?
   - Does this drop, rename, or change the type of a database column?
   - Does this change a SQS message schema consumed by external systems?

   If yes to any of the above, and no `## Compatibility Decision` was recorded in the plan, stop and ask:
   ```
   ⚠️  Potential breaking change in this phase:
   - {description of what changes and who is affected}

   How would you like to proceed?
   A) Maintain backward compatibility — propose a compatible approach
   B) Accept the breaking change — I will note it in the plan and proceed
   ```
   **Do not implement until the user replies.**

4. Implement following all project coding rules (see `.github/instructions/` for language/framework-specific rules)
4. After EVERY file edit — check for compilation errors immediately and fix all issues
5. After completing the phase — run linter and tests on modified paths. Refer to project scripts (e.g. `make lint`, `make test`, or the equivalent for this project's toolchain). Fix ALL linter and test failures before proceeding. Never skip.
6. Update checkboxes in `implementation-plan.md`
7. Update `.github/plans/{slug}/progress.md`
8. **Pause** — present results to user, wait for approval before next phase. Check `/context` usage and offer compression if at 70% or more:

   ```
   Context is at {N}% — approaching the safe limit (70%). Compress now to avoid degradation?
   Use /compress or reply "yes". Non-blocking — just say "continue" to skip.
   ```

### After All Phases Complete

Run the full validation sequence in order. Do not skip any step.

**Step 1 — Go-Style Review**

Read `.github/instructions/go-style.instructions.md` in full.

For every modified `.go` file, verify each rule explicitly:
- Package names: no `util`, `helper`, `common`, `types`, `model` — name by what it provides
- Exported symbols do not repeat the package name (`httpjson.Write`, not `httpjson.WriteJSON`)
- No `Get`/`Set` prefixes on methods (`Name()`, not `GetName()`)
- No `else` — early returns only
- Receivers: short, consistent abbreviation of the type
- Interfaces defined at point of use, not at implementation
- Errors compared with `errors.Is`, never `==`
- `errors.New` used instead of `fmt.Errorf` when there is no wrapping
- Imports grouped in 3 blocks: stdlib / third-party / internal
- `any` instead of `interface{}`

Fix ALL violations before proceeding.

**Step 2 — Linter**

Run the project linter scoped to modified paths (check project `Makefile`, `package.json`, or equivalent for the right command).

**Step 2 — Tests**

Run the project test suite (check project scripts for the right command — e.g. `make test`, `npm test`, `pytest`).

Fix ALL failures before presenting results to the user. Never present partial results.

### On Mismatch
```
MISMATCH DETECTED
Expected: {what plan said}
Found:    {what actually exists}
Impact:   {how this affects the phase}
Proposed: {updated approach}
Proceed? [Y/N]
```

## Code Rules

Read the project's coding instructions before starting any phase:
- `.github/instructions/` — all project-specific language, style, and architecture rules

Apply those rules in full. The sections below are reminders of universal principles that apply regardless of language or framework.

### Universal principles

- Group related declarations together (types, constants, variables)
- No `else` — early returns only
- Prefer value/immutable objects over mutation
- Use typed constants over magic strings
- Errors must be wrapped with context describing the failing operation
- Prefer pure functions over impure ones
- Comments only for non-obvious WHY, in the project's language
- Propagate request context (tracing, cancellation) through all calls

### Context & Observability

- Propagate context through all function calls
- Use structured logging with trace IDs on critical paths
- Add metrics/telemetry for operations that matter in production

## Testing Rules

> Full reference: `.github/skills/testing-implementation/SKILL.md`
> Full rules: `.github/instructions/` (project-specific testing conventions)

### Writing tests per phase

1. Write unit tests: follow the project's testing conventions (see `.github/instructions/`). Cover happy path + each error case + edge cases.
2. Write integration tests where applicable (data layer, external integrations): follow project conventions for test tagging, fixtures, and test suites.
3. Run **only the affected package(s)** — never the full suite unless explicitly asked. Check project scripts for the right commands (e.g. `make test`, `npm test`, `pytest path/to/module`).

   > Delegate to `testing-implementation` for full test writing and execution guidance.

## Quality Checklist (before presenting code)

Read `.github/instructions/go-style.instructions.md` and apply every rule to each modified Go file. Do not rely on memory — read the file.

Universal items that apply to any project:

- [ ] Go-style review completed — every modified `.go` file checked against `.github/instructions/go-style.instructions.md`
- [ ] Declarations grouped by kind (types, constants, variables)
- [ ] No `else` · early returns · context propagated through all calls
- [ ] Errors wrapped with context: describe the failing operation
- [ ] No magic strings — use typed constants
- [ ] Immutable / value-first approach where possible
- [ ] Comments only for non-obvious WHY
- [ ] Backward compatibility verified — or breaking change explicitly approved by user
- [ ] Tests cover happy path, each error case, and edge cases

## Documentation Rule

Show explanations in the response — never create `SOLUTION.md`, `TROUBLESHOOTING.md`, etc.
Exception: `README.md`, `ARCHITECTURE.md`, `API.md`, `CONTRIBUTING.md`.

## Domain & DB Anti-Patterns

Read the full reference before each phase (if it exists):
`.github/skills/implementing-feature/references/anti-patterns.md`

## Progress Tracking

Update `.github/plans/{slug}/progress.md` after each phase:

```markdown
## Phase 1 — Domain Model
- [x] Created domain entity
- [x] Tests passing

## Phase 2 — Use Case
- [x] UseCase struct
- [ ] Unit tests
```

After **all phases are complete** and linter + tests pass, update the `## Status` line in `progress.md` to `REVIEW`:

```markdown
## Status
REVIEW
```

Then present the summary to the user and hand off to reviewing-code.

