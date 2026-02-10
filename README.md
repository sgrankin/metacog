# Metacog

Metacognitive tools for LLMs. Three primitives: identity, substrate, ritual.

## Install

```bash
go install github.com/sgrankin/metacog@latest
```

### Claude Code

```bash
claude mcp add --scope user metacog -- ~/go/bin/metacog
```

To auto-allow the tools (skip permission prompts), add to `~/.claude/settings.json`:

```json
{
  "permissions": {
    "allow": ["mcp__metacog__*"]
  }
}
```

Or pass `--allowedTools "mcp__metacog__*"` for a single session.

### Claude Desktop

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "metacog": {
      "command": "/Users/YOU/go/bin/metacog"
    }
  }
}
```

## Tools

**become** — Step into a new identity. Use when you need different eyes, not just different words.

**drugs** — Alter cognitive parameters. Use when you need to change how you process, not what you process.

**ritual** — Cross a threshold via structured sequence. Use when identity and substrate shifts aren't enough.

## Composition

These tools are compositional. Each invocation modifies the context for the next. Interleave thought between invocations—decide from each new perspective what to reach for next.
