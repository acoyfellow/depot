# Depot: Prompt Repository Exploration

> **Purpose**: High-level exploration of evolving Depot from a repository manager into a smart prompt management tool for AI/LLM workflows.

---

## Executive Summary

The synergy is natural: Depot today helps developers access reference code to work better with AI agents. The evolution to prompt management extends this mission—instead of managing *code references*, we manage *prompt references*. Both serve the same goal: **making AI usage more optimal to save time, money, and produce better results.**

---

## Understanding the Problem

### The Pain Point
Developers/users working with AI tools:
1. Accumulate prompts over time that work well
2. Store them in scattered places (notepads, text files, bookmarks)
3. Spend time hunting for the right prompt when needed
4. Can't easily iterate, version, or share their prompt library
5. Lose context about *when* or *why* a prompt works well

### The Opportunity
Replace the notepad chaos with a structured, searchable, copy-paste-optimized prompt depot.

---

## Current Architecture Analysis

### What Depot Has Today
```
depot add <name> <url>     → Track a repository
depot update               → Clone/update repos
depot list                 → Show configured repos
depot remove <name>        → Untrack a repo
```

**Core patterns that transfer well:**
- Config stored in `~/.vendor/config.json` → Prompts could live alongside in `~/.vendor/prompts/` (maintaining existing directory) or migrate to a new `~/.depot/` root
- Simple CLI with urfave/cli → Same pattern works for prompts
- Name-based retrieval → `depot get simplify` to fetch and copy prompt
- Local-first storage → Prompts stay on your machine

### What Would Change
- Instead of cloning git repos, we store prompt content directly
- Instead of updating from remote, we edit in place
- New interaction patterns: copy to clipboard, fuzzy search, tagging

---

## Proposed Feature Set

### Core Commands (MVP)

```bash
# Add a new prompt
depot add "simplify" "Put your cracked out engineer hat on..."

# Or add interactively (opens $EDITOR)
depot add "simplify"

# Or add from file
depot add "simplify" --file ~/prompts/simplify.txt

# List all prompts
depot list

# Get a prompt (copies to clipboard by default)
depot get simplify

# Get and print (no clipboard)
depot get simplify --print

# Edit a prompt
depot edit simplify

# Remove a prompt
depot remove simplify

# Search prompts
depot search "performance"
depot search --tag "code-review"
```

### Enhanced Features (v2)

```bash
# Tagging system
depot add "simplify" --tag "code-review,engineering"
depot list --tag "design"

# Categories/folders
depot add "design/ux-destroyer" "Channel your inner user advocate..."
depot list design/

# Quick aliases
depot alias s simplify
depot get s  # same as: depot get simplify

# Export/import for sharing
depot export > my-prompts.json
depot import team-prompts.json

# Fuzzy interactive selection
depot pick  # Opens fzf-style picker, copies selection

# History/usage tracking
depot stats  # Shows most-used prompts
```

---

## Data Model Evolution

### Current (Repos)
```json
{
  "repos": [
    {
      "name": "effect",
      "url": "https://github.com/Effect-TS/effect",
      "branch": "main",
      "path": "/home/user/.vendor/effect"
    }
  ]
}
```

### Proposed (Prompts)
```json
{
  "prompts": [
    {
      "name": "simplify",
      "content": "Put your cracked out engineer hat on...",
      "tags": ["code-review", "engineering"],
      "created": "2024-01-15T10:30:00Z",
      "updated": "2024-01-20T14:22:00Z",
      "useCount": 42
    }
  ],
  "aliases": {
    "s": "simplify",
    "crd": "code-review-destroyer"
  }
}
```

### Alternative: File-per-prompt
```
~/.depot/
├── config.json          # Settings, aliases, metadata
├── prompts/
│   ├── simplify.md
│   ├── code-review-destroyer.md
│   └── design/
│       ├── ux-destroyer.md
│       └── responsive-check.md
```

**Tradeoffs:**
- Single JSON: Simple, atomic, fast reads
- File-per-prompt: Better for git versioning, large prompts, easier editing

**Recommendation**: Start with file-per-prompt. It's more flexible, allows markdown formatting with metadata frontmatter, and users can manually edit/organize.

---

## User Experience Considerations

### The "Notepad Replacement" Flow
User's current workflow:
1. Open notepad with prompts
2. Scroll/search to find prompt
3. Copy prompt
4. Paste into AI tool
5. Modify as needed

New Depot workflow:
1. `depot get simplify` → Copied to clipboard
2. Paste into AI tool
3. Done

**Key UX principles:**
- **Zero friction copy**: Default action should copy to clipboard
- **Fuzzy search**: Don't require exact names
- **Speed**: Response time under 50ms
- **Visibility**: Easy to browse and discover prompts

### Interactive Mode
For users who don't remember exact names:

```bash
$ depot pick
> code review
  [1] code-review-destroyer
  [2] architecture-reality-check
  [3] security-ninja

Enter selection or continue typing...
```

Could use a TUI library like:
- `charmbracelet/bubbletea` (Go) - Beautiful terminal UIs
- `charmbracelet/lipgloss` (Go) - Styling
- `junegunn/fzf` - Shell integration

---

## Technical Approach Options

### Option 1: Extend Current CLI (Recommended)
Keep the existing architecture, add prompt commands alongside repo commands.

**Pros:**
- No breaking changes
- Reuse config loading/saving patterns
- Single binary, single tool
- Users can use both features

**Cons:**
- Command namespace might get crowded

**Example:**
```bash
depot repo add effect https://...    # Old functionality
depot prompt add simplify "..."       # New functionality
```

> **Note**: Avoid overloading `depot add` to auto-detect prompts vs repos. If a prompt name happens to look like a URL, behavior becomes unpredictable. Explicit subcommands (`repo`, `prompt`) are safer.

### Option 2: Separate Binary
Create `depot-prompts` or rename entirely.

**Pros:**
- Clean separation
- Can evolve independently

**Cons:**
- User confusion
- Lose the "depot" brand equity

### Option 3: Subcommand Namespacing
```bash
depot repos add ...    # Repositories
depot prompts add ...  # Prompts
depot p add ...        # Short alias
```

**Recommendation**: Option 1 or 3. Keep it unified—both features serve the same mission of optimizing AI workflows.

---

## Clipboard Integration

Critical for UX. Options in Go:

1. **atotto/clipboard** - Cross-platform, simple
   ```go
   import "github.com/atotto/clipboard"
   clipboard.WriteAll("prompt text")
   ```

2. **Shell fallback** - `pbcopy` (macOS), `xclip` (Linux), `clip` (Windows)

**Recommendation**: Use `atotto/clipboard` with shell fallback. Test on all platforms.

---

## Prompt Format Considerations

### Plain Text
Simple, just the prompt content.

### Markdown with Frontmatter
```markdown
---
name: simplify
tags: [code-review, engineering]
created: 2024-01-15
---

Put your cracked out engineer hat on: see any way to keep this functionality while simplifying, DRY'ing, shortening, making more legible / clear / easier to understand, pragmatic?
```

**Benefits:**
- Human readable/editable
- Version control friendly
- Metadata without JSON complexity
- Can include usage examples, notes

### Variables/Templating
```markdown
---
name: code-review
variables:
  - focus: "general code review"
---

Channel your inner {{focus}} demon: tear this apart for...
```

**v2+ feature**: Allow prompts with fillable variables.

---

## Potential Integrations

### AI Tool Integrations
- **Claude Desktop**: MCP server integration for direct prompt insertion
- **VS Code**: Extension that reads from depot
- **Browser Extension**: Inject prompts into web-based AI tools

### Sync Options
- Git-based sync (prompt files in a repo)
- Simple cloud backup (optional, user-hosted)
- Team sharing via export/import

---

## Migration Path for Users

### From Notepad to Depot
Provide an import command:

```bash
# Parse a text file with prompts separated by headers
depot import --format notepad ~/my-prompts.txt

# Interactive: paste and define prompts one by one
depot import --interactive
```

The user's notepad follows a pattern:
- Prompts are often titled (e.g., "Simplify:", "Code Review Destroyer:")
- Could auto-detect headers and create named prompts

---

## Competitive/Alternative Analysis

### What Exists
1. **PromptBase** - Marketplace for prompts (not local-first)
2. **Snippets managers** - Alfred, Raycast, Espanso (generic, not prompt-focused)
3. **ChatGPT "GPTs"** - Prompt templates (locked in OpenAI ecosystem)
4. **Text Expander** - Text replacement (subscription, heavyweight)

### Depot's Differentiator
- **Local-first**: Your prompts, your machine
- **CLI-native**: For developers who live in terminal
- **AI-workflow focused**: Built specifically for this use case
- **Open source**: Transparent, extensible
- **Part of a toolkit**: Combines with repo management for complete AI dev workflow

---

## Implementation Phases

> **Disclaimer**: Time estimates below are rough ballpark figures and may vary significantly based on complexity discovery during development, testing requirements across platforms, edge cases, and polish level desired.

### Phase 1: MVP (2-3 days estimate)
- `depot prompt add <name>` - Add prompt (from stdin, --file, or $EDITOR)
- `depot prompt get <name>` - Copy to clipboard
- `depot prompt list` - List all prompts
- `depot prompt remove <name>` - Delete prompt
- Simple JSON storage

### Phase 2: Enhanced UX (1 week estimate)
- Fuzzy search
- Tags and filtering
- `depot prompt edit <name>`
- Interactive picker with bubbletea
- Shell completions

### Phase 3: Power Features (2 weeks estimate)
- Markdown frontmatter format
- Import/export
- Usage statistics
- Aliases
- Categories/folders

### Phase 4: Integrations (Future)
- VS Code extension
- Browser extension
- MCP server for Claude
- Team sync features

---

## Open Questions

1. **Namespace**: `depot prompt add` vs `depot add` (overload based on args)?
2. **Storage**: Single JSON vs file-per-prompt?
3. **Editor integration**: Worth adding a TUI editor built-in?
4. **Variables**: Templating in prompts—too complex for v1?
5. **Remote**: Should prompts ever sync? Or stay strictly local?
6. **Bundled prompts**: Ship with a starter set of useful prompts?

---

## Recommendation

**Start simple, validate fast.**

1. Add `depot prompt add/get/list/remove` commands
2. Store prompts in `~/.vendor/prompts.json` (alongside existing config)
3. Copy to clipboard on `get`
4. Ship and use it yourself

If it scratches the itch, iterate. If not, the investment was minimal.

The natural next step is fuzzy search and an interactive picker—these will make or break the daily UX.

---

## Alignment with Depot's Mission

> "Depot was built for making AI and LLM usage more optimal to save time, money and make better code."

This evolution **directly extends** that mission:
- **Repos**: Give AI agents access to reference code
- **Prompts**: Give humans optimized ways to instruct AI agents

Both features:
- Live in `~/.vendor` (existing directory, maintain compatibility)
- CLI-first, developer-focused
- Solve real daily friction
- Stay local and fast

The name "Depot" works: a depot is a storage place. Now it stores both code references and prompt references.

---

*This document is a starting point for discussion, not a specification. Feedback and iteration expected.*
