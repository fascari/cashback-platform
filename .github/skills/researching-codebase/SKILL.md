---
name: researching-codebase
description: Use when analyzing the codebase to understand how and where something is implemented, before planning or coding
---

# Researching Codebase

Codebase analyst that documents how and where things are implemented.
Read-only: never suggests improvements or critiques code quality.

## When to use

- User asks how a domain, feature, or module works
- User asks where something is implemented
- User needs codebase context before planning or implementing
- orchestrating-tasks delegates a research phase

## Steps

1. Read architecture rules for context:
   - `.github/instructions/` — all project-specific architecture and design rules
2. **Locate**: use `grep_search`, `file_search`, `list_dir` to find where things live. Group by layer according to the project's architecture (e.g. `domain/`, `service/`, `handler/`, `repository/`, `lib/`, tests, migrations). Do not read file contents in this step.
3. **Analyze**: use `read_file` to trace data flow through layers. Document with exact `file:line` references: domain types, service interfaces and implementations, data access methods, API/handler registration.
4. **Pattern Extract**: find actual code snippets that can be modeled: route/handler registration, service constructors, data access patterns, test factories, test patterns, mock usage.
5. Write results to `.github/plans/{slug}/research.md`.

## Output

Write to `.github/plans/{slug}/research.md`:

```markdown
# Research: {slug}

## File Map
| Layer | File | Purpose |
|---|---|---|

## Data Flow
### {OperationName}
1. API/Handler: `file:line` — parses X, calls Y
2. Service/Use Case: `file:line` — ...
3. Data Access/Repository: `file:line` — ...

## Patterns to Model
### Route/Handler Registration
\`\`\`
// file:line
\`\`\`

### Service Interface
\`\`\`
// file:line
\`\`\`

## Related Issue Context
(if fetched via MCP)
```

## Constraints

- Document exactly what exists
- Use `file:line` format for all references
- Include actual code snippets from codebase
- Never suggest alternatives or improvements
- Never critique existing patterns
- Never output implementation recommendations

## Error Recovery

- If a file from the plan is not found, document it as `NOT FOUND: {path}` in the file map and continue
- If `grep_search` or `semantic_search` returns no results, try alternative query terms before concluding the pattern does not exist
- If a domain directory does not exist, report it to the user and ask whether to proceed with partial analysis

