---
name: creating-pull-request
description: Use when opening a GitHub pull request for a branch that has commits ready for review
---

# Creating Pull Request

Analyzes the current branch, gathers context from commits and changed files,
generates a complete Pull Request description following the project template,
determines the correct assignee, labels, and target branch, then opens the PR
via GitHub CLI after explicit user approval.

## When to use

- User asks to open or create a PR
- User says "open a pull request", "create PR for this branch"
- User types /create_pr
- Branch is ready for review with commits pushed

## Steps

### Prerequisites

- GitHub CLI (`gh`) must be authenticated: `gh auth status`
- You must be on a feature/bugfix/hotfix branch (not `main` or `develop`)



### Step 1 — Gather Branch Context

```bash
git --no-pager branch --show-current
git --no-pager log main..HEAD --oneline
git --no-pager log develop..HEAD --oneline
git --no-pager diff main --stat
git --no-pager diff develop --stat
```

Parse the branch name to extract:
- **Type**: `feature` | `bugfix` | `hotfix`
- **Base branch**: see table below

| Branch prefix | Base branch | Default label |
|---|---|---|
| `feature/` | `develop` or `main` | `feature` |
| `bugfix/` | `develop` or `main` | `bug` |
| `hotfix/` | `main` | `bug`, `hotfix` |

If base branch is ambiguous, check which remote branch has fewer divergent commits:

```bash
git --no-pager rev-list --count HEAD ^origin/develop 2>/dev/null
git --no-pager rev-list --count HEAD ^origin/main
```

### Step 2 — Analyze Commits and Changed Files

```bash
git --no-pager log develop..HEAD --format="%s%n%b" 2>/dev/null || git --no-pager log main..HEAD --format="%s%n%b"
git --no-pager diff develop --name-only 2>/dev/null || git --no-pager diff main --name-only
```

Read the most relevant changed files to understand:
- What domains/features were touched
- Whether new endpoints were added (→ API docs / permissions changed?)
- Whether migrations exist (→ note in description)
- Whether tests were added/updated

### Step 3 — Determine PR Metadata

**Title**: a brief, imperative description of the change.
```
{brief description of the change}
```
Example: `Add get user endpoint`

**Assignee**: current git user
```bash
git --no-pager config user.email
gh api user --jq '.login'
```

**Labels** — choose ALL that apply (labels can be combined):

| Condition | Label(s) |
|---|---|
| Branch is `feature/` | `feature` |
| Branch is `bugfix/` | `bug` |
| Branch is `hotfix/` | `bug` + `hotfix` |
| New feature or improvement (not a bug) | `enhancement` |
| Auto-backport from main → develop | `backport` |
| Part of a versioned release | `release` |
| Only docs/config changed | `documentation` |
| Dependency file updated | `dependencies` |

Labels are **not mutually exclusive** — a `feature` PR that also improves existing behaviour should have both `feature` and `enhancement`.

**Reviewers**: do NOT set automatically — user will assign.

### Step 4 — Generate PR Body

Use EXACTLY the project template structure from `.github/pull_request_template.md`.

Fill in each section:

**Description section**: write 2–5 sentences explaining:
1. **Context**: what problem this solves or what feature this adds
2. **What changed**: the key technical changes (domains, use cases, handlers, migrations)
3. **Impact**: any side effects, breaking changes, or things reviewers should pay special attention to

Use the commit messages and changed file list as primary source.

**Author checklist**: pre-check items that are clearly satisfied based on analysis:
- ✅ Commit logs reviewed → always check if commits follow convention
- ✅ PR matches title/description → always check
- ✅ Test coverage → check only if test files were added/modified
- ✅ Documentation → check only if docs were updated
- ✅ API spec validated → check only if API spec files were modified
- ✅ Performance → check only if no unnecessary loops identified in diff
- ✅ Code clarity → check only if code follows project standards
- ✅ Tests assert correctly → check only if tests were written

**References section**: keep as-is from template.

### Step 5 — Present Plan for Approval

Present the full PR plan before any action:

```
PR Plan:

  Title:       Add get user endpoint
  Base branch: main
  Assignee:    @me
  Labels:      feature, enhancement
  Draft:       No
  Push:        git push

  ── Body Preview ──────────────────────────────────────────

  ## Description
  This PR adds the get user endpoint that...
  ...

  ## Author' checklist
  * [x] Did I review my commit logs...
  ...

  ──────────────────────────────────────────────────────────

  Command to be executed:
  gh pr create \
    --title "Add get user endpoint" \
    --base main \
    --assignee "@me" \
    --label "feature" \
    --body-file /tmp/pr-body.md

Open this PR? [Y/N] — or type 'draft' to open as Draft PR
```

**WAIT for explicit approval**: "yes", "ok", "y", "proceed", or "draft".

### Step 6 — Write Body to Temp File and Open PR

**ALWAYS use the `create_file` tool** to write the PR body — never use shell heredoc (`cat > file << 'EOF'`) or inline Python.
Shell heredocs lock the terminal and corrupt multi-line content in zsh. The `create_file` tool is safe, reliable, and produces correct Markdown.

#### Body formatting rules (enforced via `create_file`)

- Every Markdown section (`## Title`) must be preceded by a **blank line**
- Every paragraph inside `## Description` must be separated by a **blank line**
- Each sentence/paragraph must be on a **single unbroken line** — no mid-sentence line wraps
- Blockquotes (`>`) must have a blank line before and after them
- Checklist items must have a blank line after the section header (`## Author' checklist`)

#### Example — correct body file content

```markdown
## Description

This PR adds the get user endpoint that retrieves a user profile by ID.

The new handler reads the user ID from the request path and delegates lookup to the service layer.

> **Note:** any relevant side-effect or cherry-pick note goes here as a single line.

## Author' checklist

* [x] Did I review my commit logs to make sure all commits are atomic and well described?
* [x] This PR does only and exactly what the title and description are saying?
```

#### Workflow

1. Push the branch before creating the PR:
   - If the user passed `--no-verify` → `git push --no-verify`
   - If the project has pre-push validation hooks (linting, API spec validation, etc.) → `git push` (no `--no-verify`) so hooks run
   - Otherwise → `git push`
2. Use `create_file` tool to write the body to `/tmp/pr-body.md`
3. Then run:

```bash
# Open the PR
PAGER=cat gh pr create \
  --title "{title}" \
  --base {base-branch} \
  --assignee "@me" \
  --label "{label1}" \
  --label "{label2}" \
  --body-file /tmp/pr-body.md
```

If user typed `draft`, add `--draft` flag.

> **Note on labels**: Labels must already exist in the GitHub repository. If a label does not exist, `gh pr create` will return an error. In that case, create the PR without the missing label and inform the user:
> ```bash
> PAGER=cat gh label create "migration" --color "#0075ca" --description "Contains DB migration"
> ```

### Step 7 — Verify and Present Result

```bash
gh pr view --web   # opens in browser (optional, ask user)
gh pr view         # shows PR summary in terminal
```

Present the PR URL and a summary to the user.

## PR Body Template

Read the project PR template via `read_file`: `.github/pull_request_template.md`

Use that template structure when generating the body. Fill in each section based on branch context,
commits, and changed files. Pre-check Author' checklist items that are clearly satisfied.

## Labels

Use `PAGER=cat gh label list` to discover available labels. Choose ALL that apply:

| Branch prefix | Default label(s) |
|---|---|
| `feature/` | `feature` |
| `bugfix/` | `bug` |
| `hotfix/` | `bug` + `hotfix` |

Labels are not mutually exclusive — combine freely (e.g., `feature` + `enhancement`).

## Anti-patterns (never do these)

| Wrong | Correct |
|---|---|
| Open PR without user approval | Always show plan and wait for `[Y/N]` |
| Hardcode reviewer usernames | Never set reviewers automatically |
| Use shell heredoc for body (`<< 'EOF'`) | Use `create_file` tool to write `/tmp/pr-body.md` |
| Use inline Python one-liner to write body | Use `create_file` tool to write `/tmp/pr-body.md` |
| Paragraphs without blank lines between them | Every paragraph separated by one blank line |
| Mid-sentence line breaks in description | Each sentence/paragraph on a single unbroken line |
| Run `gh` commands without `PAGER=cat` | Always prefix `gh` with `PAGER=cat` |
| Guess labels that don't exist | Check label existence, create if missing |
| Skip checklist items silently | Mark unchecked items with `[ ]` + note why |
| Run `git push` for pure docs/config changes when pre-push hooks don't apply | `git push --no-verify` — hooks are not relevant |
| Use `--no-verify` when project validation hooks must run | Always `git push` (no `--no-verify`) when hooks are mandatory |
| Add `Co-authored-by: Copilot` to the PR body | The `Co-authored-by` trailer belongs only in **git commit messages** — never in the PR description or any GitHub comment |





