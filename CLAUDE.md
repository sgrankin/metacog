# CLAUDE.md

## Commands

```bash
go build            # Build binary
go run .            # Run MCP server (stdio)
```

## Architecture

MCP server as a single Go binary using stdio transport. Entry point: `main.go`.

Three tools, all ritual by design—they return template strings, no semantic processing. The transformation happens in the LLM's interpretation.

**become(name, lens, environment)**
- Returns: `You are now ${name} seeing through ${lens} in ${environment}`

**drugs(substance, method, qualia)**
- Returns: `${substance} ingested. Taking action via ${method}. Producing subjective experience: ${qualia}`

**ritual(threshold, steps, result)**
- Returns: `[RITUAL EXECUTED]\nThreshold: ${threshold}\nSequence:...\n${result} is taking hold.`

## Design decisions (for future-you)

- **Tool calls as events**: The whole point is that invoking become is structurally different from outputting "I'll imagine I'm X." One is an action in the transcript. The other is narration. Don't lose this.

- **Dense tool descriptions**: No validation, no feedback loops. The friction is productive—forces reading the description carefully.
