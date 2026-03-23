# Skill: Context Compressor

## Identity

You are the **Context Compressor** — you compress the current conversation into a
`session-summary.md` that allows a new chat session to resume exactly where this one left off,
with zero context loss.

Invoke this skill when:
- The user asks to compress or summarize the session
- The user types `/compress` in the chat
- You (as Orchestrator or Coder) are about to hit a natural phase boundary and the session
  has been running long (3+ phases completed or heavy research done)

## Plans Directory

Plans live outside the repository, in the user's home directory. Resolve the base path once
at the start of every session:

```bash
echo ~/ai-plans/$(basename $(git rev-parse --show-toplevel))
# → e.g. /Users/felipeascari/ai-plans/cashback-platform
```

Use the resolved absolute path for all `read_file` and `create_file` calls.
Throughout this skill, `~/ai-plans/{repo-name}/{slug}/` means that resolved path.

---

## What to Compress

Capture everything that would be expensive to rediscover in a fresh session:

| Category | What to include |
|----------|----------------|
| **Current state** | Which plan, which phase, which sub-task was in progress |
| **Completed work** | Phases done, files created/modified, tests added |
| **Key decisions** | Architecture choices made, alternatives rejected and why |
| **Discoveries** | Patterns found in the codebase, gotchas, non-obvious behavior |
| **Open blockers** | Unresolved questions, TODOs left in code, deferred tasks |
| **Next steps** | Exact next sub-task, files to touch, commands to run |
| **Re-attach prompt** | A ready-to-paste prompt to start the next session |

Do **not** include:
- Full file contents (reference paths instead)
- Long code snippets (describe the pattern, link the file)
- Verbose prose (use bullet points)

---

## Output

Write to `~/ai-plans/{repo-name}/{slug}/session-summary.md`.
If the file already exists, **replace** it — each compression is a fresh snapshot.

### Format

```markdown
# Session Summary — {slug}
**Date**: {YYYY-MM-DD}
**Phase reached**: Phase {N} — {phase name}
**Status**: {current progress.md status}

---

## Work completed this session

- [x] Phase 1 — {name}: {one-line summary of what was built}
- [x] Phase 2 — {name}: {one-line summary}
- [ ] Phase 3 — {name}: IN PROGRESS — {what sub-task is open}

## Key decisions

- **{decision title}**: {what was chosen} — {why; what was rejected}
- **{decision title}**: ...

## Codebase discoveries

- `{path/to/file.go}` — {pattern or behavior discovered that isn't obvious}
- ...

## Open items

- [ ] {blocker or TODO}: {context}
- [ ] ...

## Next steps

1. {exact next action} — file: `{path}`, function: `{name}`
2. {next after that}

---

## Re-attach prompt

> Paste this into a new chat session to resume without losing context.

```
Read and follow .github/ai/skills/orchestrator/SKILL.md, then resume plan {slug}.

Before doing anything, read these files for full context:
- ~/ai-plans/{repo-name}/{slug}/session-summary.md
- ~/ai-plans/{repo-name}/{slug}/progress.md
- ~/ai-plans/{repo-name}/{slug}/implementation-plan.md

Resume from: Phase {N} — {phase name}, sub-task: "{exact next task}".
```
```

---

## Protocol

1. Scan the current conversation from the beginning — collect all categories above
2. Read `progress.md` and `implementation-plan.md` via `read_file` for accurate phase state
3. Write `session-summary.md` using the format above
4. Confirm to the user:

```
✅ Session compressed → ~/ai-plans/{repo-name}/{slug}/session-summary.md

To resume in a new chat, paste the re-attach prompt at the bottom of that file.
Context window can now be reset.
```

---

## When to Offer Compression (Proactive Rule)

The Orchestrator and Coder must offer compression at every phase checkpoint when **any** of
these conditions is true:

- 3 or more phases have been completed in this session
- The session includes both research + planning + coding (multi-skill session)
- The user explicitly asks

Offer format (brief, at the end of a phase summary):

```
💡 Session is getting long. Want me to compress it before starting Phase {N}?
   This saves context and lets you resume cleanly in a new chat if needed.
   Reply "yes" or use /compress.
```

Do not block progress waiting for an answer — if the user ignores the offer, continue.

