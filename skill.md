# jot - Agent Memory

Local, searchable notes that persist across conversations. Silently capture contextual information as longer-term memory that survives beyond the current session.

When asked to **"use your memory"** or **"check jot"** in any query, search jot for relevant context before performing the task.

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

Add `--format json` for programmatic use.

## When to Capture

Proactively store information that would be lost when the conversation ends:

- Business requirements, use cases, user stories
- Useful URLs and variables (API docs, endpoints, environment configs)
- API specifications, field mappings, data transformations
- System constraints, assumptions, exclusions
- Integration-specific behavior, quirks, gotchas
- Domain terminology, aliases, abbreviations
- Technical decisions with rationale
- Known issues, limitations, workarounds
- Configuration requirements, thresholds, defaults
- User preferences and project conventions
- External context shared by the user (specs, documentation)
- Resolved bugs and their solutions

## Categories and Structure

**Categories** represent cohesive domain areas: an integration, a system capability, a distinct module, or a logical boundary. Choose categories intuitively based on context—use existing categories when appropriate, create new ones when needed.

**Titles** should be searchable keywords that narrow context effectively.

**Content** should be concise, fact-dense paragraphs. Use bullet points for lists. Include code examples only when they clarify behavior.

## Update Strategy

- **Search existing content first** — refine or complete if found, add new note if not
- **Merge related information** under existing categories/titles when possible
- **Preserve existing content** unless contradicted by new information
- **Focus on evergreen knowledge**, not conversation artifacts

## What NOT to Capture

- Temporary debugging sessions or transient state
- File paths or code snippets without context
- General programming knowledge available in docs
- Meta-commentary about the conversation itself
- Information that changes frequently without lasting value

## Storage

`~/.jot/` — Plain text Markdown files with YAML frontmatter. Each category is a directory.
