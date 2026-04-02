# Skills

AI-assisted workflow skills for any project. Each subdirectory contains a `SKILL.md`
with the full instructions for that skill.

## Entry Point

Use the **orchestrating-tasks** skill as the single entry point for any task involving codebase
analysis or changes. It detects complexity, delegates to the right skills, and
manages plan state.

```
/skill:orchestrating-tasks
```

## Skill Catalog

| Skill | Purpose | Invocation |
|---|---|---|
| orchestrating-tasks | Detect complexity, delegate, manage plan state | `/skill:orchestrating-tasks` |
| researching-codebase | Read-only codebase analysis with `file:line` references | `/skill:researching-codebase` |
| planning-implementation | Translate research into phased implementation plans | `/skill:planning-implementation` |
| implementing-feature | Implement plans phase by phase with linter and test validation | `/skill:implementing-feature` |
| reviewing-code | Code review against all project rules | `/skill:reviewing-code` |
| committing-changes | Group changes into logical commits following conventions | `/skill:committing-changes` |
| creating-pull-request | Open a GitHub PR with title, labels, and description | `/skill:creating-pull-request` |
| testing-implementation | Write and execute unit and integration tests | `/skill:testing-implementation` |
| compressing-context | Compress the current session into `session-summary.md` | `/skill:compressing-context` |
| resuming-context | Restore context from a compressed session and resume | `/skill:resuming-context` |
| sanitizing-text | Remove AI-sounding language and formatting issues | `/skill:sanitizing-text` |
| creating-adr | Create Architecture Decision Records (ADRs) | `/skill:creating-adr` |

## Standard Workflow

```
User Request → orchestrating-tasks → researching-codebase → planning-implementation → implementing-feature → reviewing-code → sanitizing-text → committing-changes
```

Full diagram: [`.github/skills/orchestrating-tasks/workflow.mmd`](orchestrating-tasks/workflow.mmd)

## Project Customization

The skills are language and framework agnostic. Project-specific rules live in:

- `.github/instructions/` — coding style, architecture, testing conventions, anti-patterns

Skills reference these files at runtime. To adapt to a new project:
1. Add or update files in `.github/instructions/` with project-specific rules
2. Update `implementing-feature/references/anti-patterns.md` with codebase-specific patterns
3. The skills will apply your rules automatically during each phase
