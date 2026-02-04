# jot - Agent Memory

Local, searchable notes that persist across conversations. Use this to remember anything that should survive beyond the current session.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/MohGanji/jot/main/install.sh | bash
```

## Commands

```bash
# Save
jot add <category> --title "..." --content "..." --tags "tag1,tag2"

# Retrieve
jot search "query"
jot list [category]
jot get <category> "pattern"

# Manage
jot update <id> --content "..."
jot delete <id>
jot categories
jot tags
```

## When to Use

**Always store information that:**
- Comes from external sources (user-provided context, fetched data, research findings)
- Represents historical knowledge (past decisions, resolved issues, learned patterns)
- Needs to be reused later (credentials, configurations, preferences)
- Would be lost when the conversation ends

**Examples of what to store:**
- User preferences and project-specific conventions
- API keys, endpoints, credentials, and configurations
- Gotchas, quirks, and workarounds discovered during work
- Decisions made and their rationale
- External context the user shares (documentation, specs, requirements)
- Debugging findings and solutions to problems
- Patterns and approaches that worked well
- Any information the user explicitly asks to remember

**When retrieving:** Before starting tasks, search for relevant stored context that might help.

## Storage

`~/.jot/` - Markdown files with FTS5 search. Immediately searchable, scales to 100K+ notes.

Add `--format json` for programmatic use.
