# GitHub Copilot Instructions

This file is the entry point for AI agents working in this repo. It provides project context, rule summaries, plan conventions, and skill references. Detailed rules are auto-injected via path-specific instructions in `.github/instructions/`; skills are in `.github/skills/`.

---

## Project Basics

> **Fill in for each project.**

- **Language / Stack**: {e.g. Go 1.26, Node.js 22, Python 3.12}
- **Entrypoints**: {e.g. `cmd/server/main.go`, `src/index.ts`}
- **Local Environment**: {e.g. `.env`, `mise`}

### Essential Commands

```sh
# Add the project's actual commands here, for example:
# mise run dev          # Start locally
# mise run test         # Run all tests
# mise run lint         # Lint the codebase
# mise run build        # Build artifacts
# mise run migrate      # Run DB migrations
# mise tasks            # List all available tasks
```

---

## Documentation Output Rules

Do not create physical `.md` files for explanations, guides, or troubleshooting.

- Show explanations in chat (in-memory) only
- Do not create files like `SOLUTION.md`, `TROUBLESHOOTING.md`, `EXPLANATION.md`

Only create `.md` files for project documentation that belongs in the repository:
- `README.md`, `ARCHITECTURE.md`, `API.md`, `CONTRIBUTING.md`

---

## Plans Directory

Plans are stored locally and symlinked into `.github/plans` for IDE visibility.

- **Path**: `.github/plans/{slug}/` (relative to project root)
- **Physical location**: `~/ai-plans/{repo-name}/{slug}/`
- **Symlink**: `.github/plans -> ~/ai-plans/{repo-name}/`

Before any plan file access, ensure the symlink exists — run setup from `.github/skills/references/plans-setup.md`.

### Plan Discovery

When no explicit slug is provided, scan `.github/plans/` and read `progress.md` in each directory to find `IN_PROGRESS` plans.

Plans are never committed. They survive branch switches, stash, rebase, and `git clean` because they live outside the repo.

---

## Project Architecture

> **Fill in for each project.** Replace with the actual directory layout and architecture notes.

```
/src/                   # Application source
/tests/                 # Test files
/docs/                  # Documentation
```

Detailed architecture and layer rules are in `.github/instructions/`.

---

## Path-Specific Instructions

Detailed coding rules are in `.github/instructions/` and are auto-injected by Copilot when editing matching files:

| Instruction file | Applies to | Summary |
|-----------------|------------|---------|
| `go-style` | `**/*.go` | Google Go Style Guide + project conventions |
| `clean-architecture` | `internal/app/**/*.go` | Layer rules, DI, dependency direction |
| `testing` | `**/*_test.go`, `**/testdata/**/*.go` | Table-driven, mockery, integration suites |
| `error-handling` | `**/*.go` | Domain errors, wrapping, no log-and-return |
| `package-design` | `**/*.go` | Package philosophy: provide not contain, dependency direction |

> Add or remove entries to match the instruction files that actually exist in this project.

---

## Skills Reference

Use the **orchestrating-tasks** skill as the single entry point for tasks involving codebase analysis or changes. Non-coding tasks use their dedicated skill directly.

Invoke skills directly via `/skill:name` in the Copilot chat.

| Task | Skill |
|------|-------|
| Any codebase analysis or research | `/skill:orchestrating-tasks` -> `/skill:researching-codebase` |
| Planning an implementation | `/skill:orchestrating-tasks` -> `/skill:planning-implementation` |
| Implementation | `/skill:orchestrating-tasks` -> `/skill:implementing-feature` |
| Code review | `/skill:reviewing-code` |
| Write and run tests | `/skill:testing-implementation` |
| Writing and executing commits | `/skill:committing-changes` |
| Open a Pull Request | `/skill:creating-pull-request` |
| Compress session to avoid context loss | `/skill:compressing-context` |
| Resume a compressed session | `/skill:resuming-context` |
| Sanitize text before publishing | `/skill:sanitizing-text` |

All skill files are under `.github/skills/{name}/SKILL.md`.

### Critical: Codebase Analysis

Any analysis that involves reading source files, tracing call chains, identifying impacted paths, or researching existing behavior must go through the orchestrating-tasks -> researching-codebase flow. Do not answer analysis questions from memory or partial context: always read the relevant files first.

---

## Observability

- Always propagate request context through all calls for trace ID support
- Use structured logging on critical paths
- Add metrics/telemetry for operations that matter in production

---

## PR and Branching Conventions

### PR Checklist

- [ ] Tests added/updated
- [ ] Documentation updated if behavior changed
- [ ] Migrations reviewed (if applicable)
- [ ] All relevant modules wired

### PR Title

```
Brief description of the change
```

### Branching Strategy

| Type | From | Merge to |
|------|------|----------|
| `feature` | `main` or `develop` | `main` or `develop` |
| `bugfix` | `main` or `develop` | `main` or `develop` |
| `hotfix` | `main` | `main` |

### Branch Naming

```
feature/short-description
bugfix/wrong-calculation
hotfix/urgent-fix
```

---

## Terminal Safety Rules

These rules apply to all skills that run terminal commands.

- Always `--no-pager` on all `git` commands (`log`, `diff`, `show`, `blame`)
- Always `PAGER=cat` on all `gh` commands
- Never use heredocs (`<<EOF`) in terminal — use `create_file` tool instead
- Pipe long output: `command | head -50`
- Never embed `\n` in a single `-m` string — use multiple `-m` flags
- If a command might hang, recommend it to the user instead of running it

---

## ⛔ HARD RULE — Never Commit or Push Without Explicit User Authorization

No `git commit`, `git add` + commit, or `git push` must ever be executed by this agent or by any skill unless the user has **explicitly and unambiguously authorized it in the current turn**. This means:

- "Go ahead and implement" is **NOT** authorization to commit
- "Looks good" is **NOT** authorization to commit
- "Proceed" without seeing a commit plan is **NOT** authorization to commit
- Completing an implementation task does **NOT** automatically authorize a commit
- The user must say "commit", "yes commit", "go ahead and commit", "y" (after seeing the commit plan), or equivalent explicit authorization

This rule applies to:

- All direct bash tool calls (`git commit`, `git add`, `git push`)
- All skill invocations (`committing-changes`, `creating-pull-request`)
- Fleet mode / sub-agent dispatches — sub-agents also cannot commit without the user's authorization

---

## Approval Protocol

Skills that call external APIs (GitHub), execute git commits, or run git push must follow this protocol:

1. Present a full preview of all actions before executing
2. Wait for written approval ("yes", "approve", "confirm", "ok", "y")
3. "Go ahead", "do it", or "create all" before seeing the preview is NOT approval — present the preview first
4. If scope changes during execution, stop, present the change, and wait for approval

---

## Commit Conventions

Full execution skill: `.github/skills/committing-changes/SKILL.md`

### Format

```
<subject>

<body>
```

### Field Rules

- `subject`: **HARD LIMIT: 50 characters or fewer.** Imperative mood, capitalized, no trailing period. Describe intent or impact, not the file changed.
- `body`: required when more than one file is changed or the intent is not obvious. Single paragraph, wrap at 72 chars.

> **Self-check before finalizing**: count every character in the subject line. If it exceeds 50, shorten it. Never submit a subject line longer than 50 chars.

### Chris Beams' Seven Rules (mandatory)

1. Separate subject from body with a blank line
2. **Limit the subject line to 50 characters**
3. Capitalize the subject line
4. Do not end the subject line with a period
5. Use imperative mood: "Add feature" not "Added" / "Adding"
6. Wrap the body at 72 characters
7. Use the body to explain **what** and **why**, not how

### Examples

```
Add user health check endpoint
```
`Add user health check endpoint` = 30 chars

```
Allow git stash to include plans dir

Plans content remains ignored via .gitignore wildcard. The .gitkeep anchor makes the folder known to git so stash and branch operations work correctly.
```
`Allow git stash to include plans dir` = 36 chars

### Anti-patterns

| Wrong | Correct |
|---|---|
| `Update SKILL.md to standardize plans directory paths and enhance session compression instructions` (97 chars) | `Standardize plans path in SKILL.md` (35 chars) |
| Subject describes the file changed | Subject describes the intent/impact |
| Subject exceeds 50 chars | Shorten: ruthlessly cut words |
