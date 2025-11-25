# depo

Manage a global directory of reference repositories on your system. Clone complex dependencies once, reference them everywhere.

## The Problem

When working with AI agents (or even just yourself), you often need access to real source code to solve complex problems. Documentation isn't always enough. You need to:

- Explore actual implementation patterns in libraries you depend on
- Find real usage examples in the source code
- Understand type definitions and API details
- Debug by reading the real code

Instead of downloading the same repositories into every project (bloating your disk), `depo` manages a single shared `~/.vendor` directory that you can reference from anywhere.

## Installation

```bash
go install github.com/acoyfellow/depo@latest
```

This installs to `~/go/bin/depo`. Verify it's in your PATH by running `depo list`.

## Usage

### Add a repository to track

```bash
depo add effect https://github.com/Effect-TS/effect
depo add svelte https://github.com/sveltejs/svelte
```

Repos are stored in `~/.vendor/` by default. You can customize the path:

```bash
depo add alchemy https://github.com/sam-goodwin/alchemy --branch main
```

### Clone or update all repos

```bash
depo update
```

Or update just one:

```bash
depo update effect
```

### List your repositories

```bash
depo list
```

Shows which repos are configured and cloned.

### Remove a repo from tracking

```bash
depo remove svelte
```

(The cloned directory remains; only the config entry is removed.)

## How to Use in Your Projects

Once you've cloned repositories, reference them in your `CLAUDE.md` or `AGENTS.md`:

```markdown
## Local Reference Sources

- Effect source: `~/.vendor/effect` - Use this to understand Effect types, patterns, and APIs
- Svelte source: `~/.vendor/svelte` - Reference for reactive patterns
```

Then tell your agent: "Check out `~/.vendor/effect/src` to understand how this library works."

Your agent can then search the actual source code for implementation patterns, API usage examples, and detailed type definitions that aren't always clear from documentation alone.
