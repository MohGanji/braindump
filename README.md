# braindump

Agent-friendly local memory. Store and search notes across conversations.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/MohGanji/braindump/main/install.sh | bash
```

## Agent Setup

To enable your AI agent to use braindump as long-term memory:

1. **Add the skill** — Copy [skill.md](skill.md) to your agent's skills directory (e.g., `.cursor/skills/braindump.md`)

2. **Enable in your agent rules** — Add this line to your `agent.md` or rules file:
   ```
   Use braindump as a long term memory.
   ```

The agent will then proactively store and retrieve context across conversations.

You can also manually trigger memory operations by prompting:
- **"use your brain"** — retrieve relevant context
- **"braindump this"** — store information for later

## Usage

```bash
braindump add api-creds --title "Stripe Key" --content "sk_test_..."
braindump search "stripe"
braindump list api-creds
braindump get api-creds "stripe"
```

## Commands

```bash
braindump add <category> --title "..." --content "..." [--tags "..."]
braindump search <query> [--in category] [--tag tag1,tag2]
braindump list [category]
braindump get <category> [pattern]
braindump update <id> --content "..." [--title "..."] [--tags "..."]
braindump delete <id>
braindump categories
braindump tags
```

Add `--format json` to any command for JSON output.

## Storage

```
~/.braindump/
├── api-creds/
│   └── stripe-key.md
├── api-quirks/
│   └── webhook-gotcha.md
└── .index/
    └── search.db
```

Files are markdown with YAML frontmatter. Search is SQLite FTS5.

## Performance

Benchmarked on Apple M3 Pro:

- **Add:** ~1.2ms per note (constant, even at 100K notes)
- **Search:** 23ms (10K notes) to 166ms (100K notes)
- **Get:** O(1) direct file access
- **Scales:** 100K+ notes, tested with realistic content (100-1000 words per note)

See [benchmark/results.md](benchmark/results.md) for detailed results.

## Releasing

Push a tag to trigger automatic builds. Tags follow the pattern `vYEAR.MO.DAY` (e.g., `v2026.02.03`):

```bash
git tag v2026.02.03
git push origin v2026.02.03
```

GitHub Actions builds binaries for macOS, Linux, and Windows.

## License

MIT
