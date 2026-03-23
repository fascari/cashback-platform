# Git Commit Conventions

## Format Options

**With ticket tracking:**
```
<ticket-id> | <subject>

<body>
```

**Without ticket tracking:**
```
<subject>

<body>
```

## Guidelines

- **ticket-id**: Use your project's ticket format (e.g., `PROJ-123`, `#42`, or omit if not using tickets)
- **subject**: Max 50 chars, imperative mood, capitalized, no trailing period. Describe the **intent or impact** — not the file changed
- **body**: Required when more than one file is changed or the intent is not obvious. Wrap at 72 characters

## Examples

Single file, obvious intent — no body needed:
```
Add consumer health check endpoint
```

Multiple files or non-obvious intent — body required:
```
Allow git stash to include plans directory

Plans content remains ignored via .gitignore wildcard. The .gitkeep
anchor makes the folder known to git so stash and branch operations
work correctly.
```

## Reference

Full workflow: `.github/ai/skills/git-committer/SKILL.md`


