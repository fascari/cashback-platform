---
description: Implement a plan from ~/ai-plans/{repo-name}/ phase by phase with verification
model: claude-sonnet-4-6
---

Read and follow `.github/ai/skills/coder/SKILL.md` exactly, then implement the plan.

Plan to implement (follow the discovery rules in the skill):
- If a ticket ID or slug was given in this message, use that.
- Otherwise, scan `~/ai-plans/{repo-name}/` for a plan with `## Status: IN_PROGRESS` in `progress.md`.
- If multiple are found, list them and ask which to use.
