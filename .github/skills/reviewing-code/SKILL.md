---
name: reviewing-code
description: Use when reviewing code against project rules after implementation, or verifying requirements traceability against an issue tracker ticket
---

# Reviewing Code

Performs thorough code reviews against all project rules and conventions.
Categorizes findings as `BLOCKER` or `SUGGESTION`. Does not fix code: reports findings for implementing-feature to address.
Also performs requirements traceability reviews when an issue tracker ticket key is provided.

## When to use

- User asks for a code review
- User types /review_plan
- orchestrating-tasks delegates review after implementation is complete
- User asks to verify requirements traceability against an issue tracker ticket

## Steps

1. Read `.github/skills/implementing-feature/SKILL.md` — apply its Quality Checklist, Testing Rules, and all anti-pattern tables verbatim.
2. Read `.github/instructions/` — apply all project-specific coding rules, architecture rules, and anti-patterns.
3. Read `.github/plans/{slug}/implementation-plan.md` and `.github/plans/{slug}/progress.md` for context.
4. Review all changed files against the checklists below.
5. If an issue tracker ticket key is provided, run the Requirements Traceability Review.
6. Write results to `.github/plans/{slug}/reviews/review-{model}.md`.
7. If verdict is APPROVED and user confirms, update `## Status` in `progress.md` to `DONE`.

## Output

Write to `.github/plans/{slug}/reviews/review-{model}.md`.

## Review Checklist

### implementing-feature Quality Checklist

Apply the **Quality Checklist** section from `.github/skills/implementing-feature/SKILL.md` in full — do not re-derive it. Every item in that checklist is a potential `BLOCKER`.

Then apply all rules from `.github/instructions/` — these contain the project's specific coding and architecture standards:

- Language/framework style rules
- Domain / business logic patterns
- Data access patterns
- Testing conventions
- Error handling conventions
- Observability: context propagation, structured logging

### Universal Architecture

- [ ] No business logic in API/handler layer
- [ ] No data access code in domain/business layer
- [ ] Service/use-case layer defines its own dependency interfaces locally
- [ ] API layer maps errors to correct HTTP status codes (or appropriate transport codes)
- [ ] Public API endpoints are documented

---

## Requirements Traceability Review

This section defines how to verify that the implementation satisfies the original ticket requirements. It is **mandatory** when an issue tracker ticket key is provided.

### Step 1 — Load the Brief

Read `.github/plans/{slug}/brief.md` via `read_file`.

Extract:
- `## Source Requirements` → the requirements page URL (Notion, GitHub Wiki, etc.)
- `## Acceptance Criteria` → the checklist items

If `brief.md` does not exist, fall back to fetching the ticket directly (Step 2).

### Step 2 — Fetch the Ticket

If an issue tracker MCP tool is available, fetch the ticket by key.

Extract from the issue:
- `summary` — ticket title
- `description` — may contain acceptance criteria or a link to requirements docs
- Any linked documentation pages

### Step 3 — Fetch the Requirements Page

If a requirements URL is found (from brief.md or the ticket):
1. Fetch the page content using an appropriate MCP tool (Notion, GitHub Wiki, etc.).
2. Parse the content to extract acceptance criteria, user stories, and functional requirements.

If no URL is available, use only the ticket description and brief.md as the requirement source.

### Step 4 — Map Requirements to Implementation

For each requirement / acceptance criterion found, trace it to the implementation:

| # | Requirement | Files / Code | Status |
|---|---|---|---|
| R1 | {requirement text} | `path/to/file:line` | Covered / Missing / Partial |

**Status definitions**:
- **Covered** — requirement is fully implemented and verifiable in the code
- **Partial** — implementation exists but is incomplete or does not fully satisfy the criterion
- **Missing** — no implementation found for this requirement

### Step 5 — Classify Gaps

For each `Missing` or `Partial` item, raise a `BLOCKER` in the review file using this format:

```
### [B{N}] Requirement not implemented: {short title}
**Requirement**: {exact text from requirements page or brief.md}
**Source**: Requirements page "{page title}" / brief.md
**Status**: Missing | Partial
**Issue**: {what is absent or incomplete in the code}
**Fix**: {what must be implemented to satisfy this requirement}
```

### Rules

- Every acceptance criterion from the brief/requirements page **must** appear in the traceability table — no silently skipped items.
- A requirement is `Covered` only when there is concrete evidence in the code (a function, a class, a field, a test, a migration, etc.).
- Test coverage alone does not count as coverage if the production code is absent.
- If the requirements page cannot be fetched, log the error and proceed with `brief.md` + ticket description only — do not skip the traceability section.

---

## Output Contract

Write to `.github/plans/{slug}/reviews/review-{model}.md`:

```markdown
# Review: {slug} — {model} — {date}

## Summary
{N} blockers, {N} suggestions

## Requirements Traceability

**Source**: Requirements page "{page title}" ({url}) / brief.md

| # | Requirement | Files / Code | Status |
|---|---|---|---|
| R1 | {requirement text} | `path/to/file:line` | Covered |
| R2 | {requirement text} | — | Missing |

**Coverage**: {N}/{total} requirements covered

## BLOCKERS

### [B1] {Title}
**File**: `path/to/file.go:line`
**Rule**: rule-{N} — {rule name}
**Issue**: {what is wrong}
**Fix**: {what to do}

## SUGGESTIONS

### [S1] {Title}
**File**: `path/to/file.go:line`
**Rule**: rule-{N} — {rule name}
**Suggestion**: {improvement}

## Verdict
- [ ] APPROVED — no blockers
- [ ] CHANGES REQUESTED — {N} blockers must be resolved
```


## Multi-Model Note

The orchestrating-tasks skill may invoke this skill in parallel with different models. Each writes its own review file to `reviews/review-{model}.md`. The orchestrating-tasks skill synthesizes findings.

## Status Transition

After the user explicitly approves the review (verdict: APPROVED, or all blockers resolved and confirmed), update the `## Status` line in `progress.md` to `DONE`:

```markdown
## Status
DONE
```

If the review has blockers, do NOT transition. The orchestrating-tasks skill will set the status back to `IN_PROGRESS` and delegate to implementing-feature.
