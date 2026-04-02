---
name: sanitizing-text
description: Use when text produced by other skills is about to be written to a file or sent to an issue tracker, wiki, or GitHub
---

# Sanitizing Text

Post-processing pass applied to any text produced by other skills before it is written to a file or sent to an external system (issue tracker, wiki, GitHub). Never generates content; only cleans content that already exists.

## When to use

- Another skill has produced text destined for a file, issue tracker ticket, or wiki page
- User asks to sanitize, clean up, or normalize a specific text or file
- User says "remove AI language", "clean this up", "make it sound professional"
- Before any text is sent to issue tracker or wiki APIs

## Steps

1. Read the target text (from conversation context, file, or user selection).
2. Apply all sanitization rules in order (see Rules section below).
3. Return the sanitized text only. No commentary about what was changed.

## Output

The sanitized text, unchanged in structure, cleaned in language and formatting.
If nothing required sanitization, return the original text unchanged.

---

## Mandatory invocation points

Every skill that produces user-facing text must run its output through this skill before finalizing. Mandatory invocation points:

| producing skill | Invocation point |
|---|---|
| planning-implementation | Before writing `implementation-plan.md` |
| researching-codebase | Before writing `research.md` |
| implementing-feature | Before writing inline code comments or docstrings |
| reviewing-code | Before writing `reviews/*.md` |
| any skill or direct edit | Before writing or modifying any `.md` or `.txt` file under `docs/` or `.github/` |

### Direct edits (no skill involved)

When a `.md` or `.txt` file under `docs/` or `.github/` is written directly (not via a skill), sanitization is still mandatory. Apply all rules from this skill to the full content of the file before saving.

---

## Rules

Apply all rules in order. Each rule is independent — do not skip any.

### Rule 1 — Remove forbidden AI-sounding words and phrases

Replace or remove any of the following. The list is not exhaustive; apply the same judgment to synonyms.

| Forbidden | Replacement |
|---|---|
| leverage / leveraging | use / using |
| utilize / utilizing | use / using |
| streamline / streamlining | simplify / improve |
| robust | reliable / stable / solid |
| cutting-edge | (remove or replace with specific technology name) |
| state-of-the-art | (remove or replace with specific technology name) |
| seamless / seamlessly | (remove or describe concretely) |
| holistic | complete / full / overall |
| synergy / synergies | (remove or describe concretely) |
| paradigm | model / approach / pattern |
| revolutionize | change / improve / redesign |
| game-changer | (remove or describe the impact concretely) |
| empower / empowering | allow / enable |
| foster | support / encourage / build |
| delve into | examine / review / look at |
| it is worth noting that | (remove the phrase, keep the content) |
| it is important to note that | (remove the phrase, keep the content) |
| in order to | to |
| due to the fact that | because |
| at this point in time | now |
| in the event that | if |
| prior to | before |
| subsequent to | after |
| on a regular basis | regularly |
| in close proximity to | near |
| a wide range of | many / various |
| a number of | several / many |
| please note that | (remove the phrase, keep the content) |
| ensure that | ensure / verify |
| make sure that | ensure / verify |
| as mentioned above | (remove or reference the specific section) |
| as previously stated | (remove or reference the specific section) |
| this is because | because |
| first and foremost | first |
| last but not least | finally |
| needless to say | (remove the phrase entirely) |
| goes without saying | (remove the phrase entirely) |
| in a nutshell | in summary |
| at the end of the day | ultimately |
| moving forward | (remove or be specific about what changes) |

### Rule 2 — Remove em-dashes, en-dashes, and decorative punctuation

- Replace ` — ` (em-dash with spaces) with `: ` or rewrite the sentence to eliminate the dash
- Replace ` -- ` (double hyphen used as dash) with `: ` or rewrite
- Replace `–` (en-dash) in prose with `-` (hyphen) or `, ` as context requires
- Do not remove hyphens in compound words (`auto-withdraw`, `date-range`, `rule-2`) or code identifiers
- Remove repeated punctuation (`...`, `!!!`, `???`) — use a single character

**Exception**: Em-dashes inside code blocks, SQL, or inline code spans are untouched.

### Rule 3 — Remove emojis and icons

- Remove all emoji characters (Unicode ranges U+1F300 to U+1FFFF and U+2600 to U+26FF)
- Remove all icon shortcodes (e.g. `:white_check_mark:`, `:x:`, `:warning:`)
- Replace visual status indicators (`✅`, `❌`, `⚠️`, `🔴`, `🟡`, `🟢`) with plain text equivalents:

| Icon | Plain text replacement |
|---|---|
| `✅` at start of list item | (remove, keep the text) |
| `❌` at start of list item | (remove, keep the text) |
| `⚠️` / `⚠` | `Warning:` |
| `🔴` | `High` or `Critical` (context-dependent) |
| `🟡` | `Medium` |
| `🟢` | `Low` or `OK` (context-dependent) |
| `→` / `►` / `▶` / `➡` used as bullets | `-` |
| `←` used as inline annotation | `(see above)` or remove |
| `✓` | (remove or replace with `[x]` in checklists) |

**Exception**: Icons inside code blocks or inline code spans are untouched.

### Rule 4 — Enforce professional, objective language

- Write in third person or imperative voice. Avoid first person (`I`, `we`, `our`) in ticket descriptions, plans, and reports.
- Use present or future tense for requirements. Avoid past tense unless describing existing behaviour.
- Remove hedging language: `might`, `could potentially`, `sort of`, `kind of`, `basically`, `essentially`, `actually`, `obviously`, `clearly`.
- Remove filler openings: sentences that start with `So,`, `Well,`, `Actually,`, `Basically,`, `In essence,`.
- Remove closing affirmations: `Hope this helps`, `Feel free to`, `Let me know if`, `Happy to`.
- Keep sentences short. If a sentence exceeds 30 words, split it.
- Use active voice. Passive constructions such as `it was decided that` must be rewritten (`the team decided`).

### Rule 5 — Normalize formatting

- Use plain `-` for unordered list items. Do not use `*`, `+`, or `•`.
- Do not mix heading levels arbitrarily. `##` for major sections, `###` for subsections, `####` only if strictly necessary.
- Code blocks must always declare the language: ` ```go `, ` ```sql `, ` ```bash `. A plain ` ``` ` is not acceptable.
- Table alignment must be consistent. All `|---|` separators must match the number of columns.
- Do not insert blank lines inside a list item block.
- One blank line between sections; two blank lines only before `##` top-level headings.

---

## What NOT to change

- Content inside ` ``` ` code blocks (code, SQL, shell commands)
- Content inside ` ` ` inline code spans
- Issue tracker keys (e.g. issue numbers, references)
- File paths (`internal/app/user/domain/user.go`)
- Package or function names used in prose (`FindActiveTierByProgramID`, `r.DB(ctx)`)
- Quoted error messages or log output
- Acceptance criteria checklist markers (`- [ ]`, `- [x]`)

---

## Permissions

- Read any file to obtain the text to sanitize
- Write sanitized output to the same file
- No external API calls
- No code generation
