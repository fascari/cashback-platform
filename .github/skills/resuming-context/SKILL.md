---
name: resuming-context
description: Use when resuming from a compressed session, pasting a re-attach prompt, or continuing work from a previous chat session
---

# Context Resumer

Reconstructs full working context from a compressed `session-summary.md` and hands off
execution to the correct skill so the session continues where it was interrupted.
Never starts implementation before presenting the restored state to the user.

## When to use

- User types `/resume` in the chat
- User pastes a re-attach prompt from a `session-summary.md`
- User says "resume", "continue from last session", or similar

## Steps

### Step 0 — Setup

1. Ensure the plans symlink exists — run setup from `.github/skills/references/plans-setup.md`

### Step 1 — Discover the Plan

Before matching against any user-provided input, **always** run these two commands in parallel:

```bash
git rev-parse --abbrev-ref HEAD
ls .github/plans/
```

> **Always use `ls .github/plans/`** (the symlink set up in Step 0). Never construct the raw path manually (e.g. `~/ai-plans/...`) — that risks typos and path mismatches.

Extract the identifier from the branch name using the following rules:

| Branch pattern | Extracted slug |
|---|---|
| `feature/add-user-endpoint` | `add-user-endpoint` |
| `bugfix/fix-login-error` | `fix-login-error` |
| `hotfix/urgent-fix` | `urgent-fix` |

Then apply the discovery rules:

| Situation | Action |
|-----------|--------|
| User provided full slug | Use it directly |
| Branch name matches a slug prefix | Use that slug, inform user |
| Exactly 1 plan with `IN_PROGRESS` in `progress.md` | Use it automatically, inform the user |
| Multiple `IN_PROGRESS` plans | List them and ask the user which to resume |
| No `IN_PROGRESS` plan found | Inform user — do not proceed without confirmation |

### Step 2 — Re-hydrate Context

Read **all three** files via `read_file` before doing anything else:

- `.github/plans/{slug}/session-summary.md`
- `.github/plans/{slug}/progress.md`
- `.github/plans/{slug}/implementation-plan.md`

If `session-summary.md` does not exist, stop and inform the user:

```
No session-summary.md found for plan "{slug}".
Run /compress in the original session first to generate it.
```

### Step 3 — Restore Key Decisions

Extract the **Key decisions** section from `session-summary.md`. Treat every item as a
hard constraint for the rest of the resumed session. Present them before delegating:

```
Key decisions restored (non-negotiable):
   - {decision 1}
   - {decision 2}
```

Forward these constraints verbatim to any skill delegated to.

### Step 4 — Identify Resume Point

Cross-reference `session-summary.md` sections (`## Work completed this session`, `## Next steps`)
with `progress.md` phase checkboxes to determine:

- Last **completed** phase (all checkboxes checked)
- Current **in-progress** phase (first phase with unchecked items)
- Exact **sub-task** to resume from (`## Next steps` in `session-summary.md`)

### Step 5 — Present Restored State

Show the user a concise restoration summary before delegating:

```
Session restored — {slug}

Resuming from: Phase {N} — {phase name}
   Sub-task: "{exact sub-task from Next steps}"

Completed phases:
   - Phase 1 — {name}
   - Phase 2 — {name}

In progress:
   - Phase {N} — {name}: {open sub-tasks}

Key decisions (restored):
   - {decision 1}
   - {decision 2}

Open items:
   - {open item 1}
   - {open item 2}

Proceeding to delegate to {SkillName}...
```

### Step 6 — Route to the Correct Skill

Based on the resume point and `progress.md` status, delegate to the appropriate skill:

| Resume point | Delegate to |
|---|---|
| Research incomplete | `researching-codebase/SKILL.md` |
| Plan needs to be created or approved | `planning-implementation/SKILL.md` via orchestrating-tasks |
| Implementation phase in progress | `implementing-feature/SKILL.md` |
| Tests need to be written or fixed | `implementing-feature/SKILL.md` (tests are part of the implementing-feature skill) |
| All phases done, awaiting review | `reviewing-code/SKILL.md` |
| Review approved, ready to commit | `committing-changes/SKILL.md` |

When delegating, pass:
1. The resolved slug
2. The exact resume sub-task
3. All key decisions as non-negotiable constraints
4. Any open items that affect the delegated skill's work

The delegated skill will read its own rules during its Step 0 — do not re-read rules here.

## Error Recovery

- If `session-summary.md` does not exist, stop and inform the user to run `/compress` first
- If `progress.md` and `session-summary.md` disagree on completed phases, trust `progress.md` (it is the authoritative state) and note the discrepancy to the user
- If `implementation-plan.md` does not exist, delegate to planning-implementation instead of implementing-feature
- If the slug no longer matches any plan directory, list available plans and ask the user to confirm

## Permissions

- Read any plan file (`session-summary.md`, `progress.md`, `implementation-plan.md`)
- Present restored state to user
- Delegate to any skill
- Forward key decisions as constraints
- Never write production code
- Never modify `progress.md` status
- Never skip the restored-state presentation
- Never start work without showing the user what was restored
