---
applyTo: "**/*.md,**/*.txt"
---

# Text Sanitization

Before saving any `.md` or `.txt` file, apply all rules from `.github/skills/sanitizing-text/SKILL.md`.

Key rules:
- Remove AI-sounding words (`leverage`, `utilize`, `robust`, `seamless`, `holistic`, etc.)
- Replace em-dashes with `:` or rewrite the sentence
- Remove emojis and icon shortcodes outside code blocks
- Write in imperative or third person; remove hedging language
- Use `-` for list items; declare language on every code block
