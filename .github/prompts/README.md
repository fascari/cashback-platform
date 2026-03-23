# Prompts

Slash commands for GitHub Copilot agent mode. Each file maps to a `/command` you can invoke directly in the chat.

## How to use

Type the command in the Copilot chat panel. Parameters can be passed inline.

```
/create_plan describe the feature
/commit
/implement_plan
/research_codebase
/review_plan
/compress
```

---

## Available commands

### `/create_plan`
**Model**: Claude Sonnet

Creates a detailed implementation plan through interactive research and iteration.

- Researches the codebase before asking questions
- Produces `~/ai-plans/{repo-name}/{slug}/implementation-plan.md`

```
/create_plan describe the feature
/create_plan                   # starts interactive mode
```

---

### `/implement_plan`
**Model**: Claude Sonnet

Implements a plan from `~/ai-plans/{repo-name}/` phase by phase, running tests and linting after each phase. Pauses for approval before proceeding.

```
/implement_plan                # picks up the active plan
```

---

### `/review_plan`
**Model**: Claude Sonnet

Reviews the implemented plan against all project rules (architecture, transaction pattern, Go style, testing, error handling).

Outputs findings to `~/ai-plans/{repo-name}/{slug}/reviews/review-{model}.md` categorized as `BLOCKER` or `SUGGESTION`.

```
/review_plan                   # code-rules review, picks up the active plan
```

---

### `/commit`
**Model**: GPT-4o

Analyzes staged and unstaged changes, groups them into logical commits, and presents a plan for approval before executing. Follows the team's commit conventions.

Full conventions: `.github/git-commit-instructions.md`
Full workflow: `.github/ai/skills/git-committer/SKILL.md`

```
/commit
```

---


### `/research_codebase`
**Model**: Claude Opus

Documents codebase behavior as-is. Produces research artifacts under `~/ai-plans/{repo-name}/{slug}/` for use by `/create_plan` or `/implement_plan`.

```
/research_codebase
```

---

---

### `/compress`
**Model**: Claude Sonnet

Compresses the current session into `session-summary.md` for resuming in a new chat without context loss.

```
/compress
```

---

## Relationship to skills

Prompts are thin entry points — they delegate to the skills under `.github/ai/skills/`. The skills contain the full logic, rules, and workflows.

| Prompt | Skill |
|--------|-------|
| `/commit` | `git-committer/SKILL.md` |
| `/create_plan` | `orchestrator/SKILL.md` → `researcher/SKILL.md` → `planner/SKILL.md` |
| `/implement_plan` | `orchestrator/SKILL.md` → `coder/SKILL.md` |
| `/review_plan` | `reviewer/SKILL.md` |
| `/research_codebase` | `researcher/SKILL.md` |
| `/compress` | `context-compressor/SKILL.md` |

