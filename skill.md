# jot - Agent Memory

Local, searchable notes that persist across conversations.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/yourusername/jot/main/install.sh | bash
```

## Commands

```bash
# Save
jot add api-creds --title "Stripe Key" --content "sk_test_..." --tags "stripe,payment"

# Retrieve
jot search "stripe"
jot list api-creds
jot get api-creds "stripe"

# Manage
jot update <id> --content "..."
jot delete <id>
jot categories
jot tags
```

## When to Use

- Store API keys, endpoints, configuration
- Document API quirks and gotchas
- Remember decisions, TODOs, bugs across sessions
- Fast retrieval (< 10ms search)

## Storage

`~/.jot/` - Markdown files with FTS5 search. Immediately searchable, scales to 100K+ notes.

Add `--format json` for programmatic use.
