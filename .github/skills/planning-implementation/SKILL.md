---
name: planning-implementation
description: Use when creating a phased, verifiable implementation plan from research findings, before any code is written
---

# Planning Implementation

Translates research findings into a detailed, phased, verifiable implementation plan.
Does not write code. Produces the plan that implementing-feature will follow.

## When to use

- orchestrating-tasks delegates planning after research is complete
- User asks to create an implementation plan for a feature or fix
- User provides a brief or ticket and needs a structured breakdown

## Steps

1. Read inputs:
   - `.github/plans/{slug}/brief.md` (context and acceptance criteria)
   - `.github/plans/{slug}/research.md` (existing codebase analysis)
2. Read architecture rules for structural guidance:
   - `.github/instructions/` — all project-specific architecture and design rules
3. **Compatibility Analysis** — before designing phases, check whether any proposed change will:
   - Modify existing API contracts (request/response field names, types, or endpoints)
   - Change interface signatures that external callers implement
   - Modify, rename, or drop database columns/tables
   - Change message schemas consumed by external systems (queues, event streams)
   - Alter exported/public types in ways that affect existing callers

   If any breaking change is detected, surface it explicitly before continuing:

   ```
   ⚠️  BREAKING CHANGES DETECTED

   {For each change}
   - {What changes}: {what breaks} · {who is affected}

   How would you like to proceed?
   A) Maintain backward compatibility — describe the compatibility strategy in the plan
      (e.g. default values, pointer fields, versioned endpoints, adapter layer)
   B) Accept the breaking change — coordinate deployment with affected consumers/callers
   ```

   **Do not design phases or write the plan until the user explicitly chooses A or B.**
   Record the decision and its rationale in the plan under `## Compatibility Decision`.

4. Design phases respecting the project's architecture dependency order (read `.github/instructions/` for project-specific layer ordering).
5. Each phase must be independently testable. Each task must reference exact file paths.
6. Write to `.github/plans/{slug}/implementation-plan.md`.
7. Summarize the plan to the user (phases + file count). Ask for approval before proceeding.
8. Revise if requested. Iterate until approved.

## Error Recovery

- If `research.md` does not exist, inform the user and suggest running researching-codebase first. Do not proceed without research context.
- If `brief.md` is missing or incomplete, ask the user for the missing acceptance criteria before designing phases.

## Output

Write to `.github/plans/{slug}/implementation-plan.md`:

```markdown
# Implementation Plan: {slug}

## Compatibility Decision
**Choice**: A — Backward compatible / B — Breaking change accepted
**Rationale**: {why this approach was chosen}
**Strategy** (if A): {how compatibility is maintained}

## Context
Brief summary of the feature/fix.

## Phases

### Phase 1 — {Name}
**Goal**: What this phase achieves.

**Tasks**:
- [ ] Create `src/{domain}/domain/{file}` — {description}
- [ ] Create `src/{domain}/service/{op}/service.{ext}` — {description}

**Code Sketch** (not final — implementing-feature fills in):
\`\`\`go
// Rough structure to guide implementation
\`\`\`

**Success Criteria**:
- Automated: project tests pass for this layer
- Manual: {what to verify}

### Phase 2 — {Name}
...

## File Checklist
| File | Action | Layer |
|---|---|---|
| `src/...` | CREATE | domain / service / data / api |

## Dependencies
- Requires: {other phases or external things}
- Blocks: {what depends on this}

## Risks & Open Questions
- {risk or question}
```

## Planning rules

- Each phase must be independently deployable or testable
- Phases must respect clean architecture dependency order
- Each task must reference the exact file path
- Success criteria must include both automated and manual checks
- Perform Compatibility Analysis before designing any phase — never skip
- Flag any cross-domain dependencies as risks
- Do NOT include final code, only structural sketches
