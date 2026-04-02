# Plans Directory Setup

Run before any skill that reads or writes plan files. Idempotent — safe to run multiple times.

```bash
PLANS_DIR=~/ai-plans/$(git rev-parse --show-toplevel | xargs basename)
[ -d "$PLANS_DIR" ] || mkdir -p "$PLANS_DIR"
[ -L .github/plans ] || ln -s "$PLANS_DIR" .github/plans
```

**What it does:**
- Plans live outside the repo at `~/ai-plans/{repo-name}/` — they survive `git clean`, stash, and branch switches
- `.github/plans` is a symlink for IDE visibility — never committed
- Always use `ls .github/plans/` (via the symlink), never construct the raw path manually
