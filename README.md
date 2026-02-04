# jot

Agent-friendly local memory. Store and search notes across conversations.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/MohGanji/jot/main/install.sh | bash
```

## Usage

```bash
jot add api-creds --title "Stripe Key" --content "sk_test_..."
jot search "stripe"
jot list api-creds
jot get api-creds "stripe"
```

## Commands

```bash
jot add <category> --title "..." --content "..." [--tags "..."]
jot search <query> [--in category] [--tag tag1,tag2]
jot list [category]
jot get <category> [pattern]
jot update <id> --content "..." [--title "..."] [--tags "..."]
jot delete <id>
jot categories
jot tags
```

Add `--format json` to any command for JSON output.

## Storage

```
~/.jot/
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
