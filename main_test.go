package main

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestMain intercepts test execution. When TEST_MAIN=metacog is set,
// it calls main() directly — the test binary becomes the program under test.
// Since go test -cover instruments this binary, main() gets coverage.
func TestMain(m *testing.M) {
	if os.Getenv("TEST_MAIN") == "metacog" {
		main()
		return
	}
	os.Exit(m.Run())
}

// selfExec runs the test binary as the metacog program with the given args.
// Coverage data from the subprocess is written to GOCOVERDIR if set.
func selfExec(t *testing.T, args ...string) *exec.Cmd {
	t.Helper()
	cmd := exec.Command(os.Args[0], args...)
	cmd.Env = append(os.Environ(), "TEST_MAIN=metacog")
	if coverDir := os.Getenv("GOCOVERDIR"); coverDir != "" {
		cmd.Env = append(cmd.Env, "GOCOVERDIR="+coverDir)
	}
	return cmd
}

func TestMainHelp(t *testing.T) {
	cmd := selfExec(t, "--help")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("exit error: %v\noutput: %s", err, out)
	}
	for _, want := range []string{"Six primitives", "claude mcp add"} {
		if !strings.Contains(string(out), want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestMainServer(t *testing.T) {
	cmd := selfExec(t)
	// Send a minimal JSON-RPC initialize + close stdin to trigger clean shutdown.
	cmd.Stdin = strings.NewReader(
		`{"jsonrpc":"2.0","method":"initialize","id":1,"params":{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"test","version":"0.1.0"}}}` + "\n",
	)
	out, err := cmd.CombinedOutput()
	// The server will exit with an error when stdin closes, which is expected.
	_ = err
	_ = out
}

func setup(t *testing.T) *mcp.ClientSession {
	t.Helper()
	ctx := context.Background()
	ct, st := mcp.NewInMemoryTransports()

	server := newServer()

	done := make(chan error, 1)
	go func() { done <- server.Run(ctx, st) }()

	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	cs, err := client.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		cs.Close()
		<-done
	})
	return cs
}

func TestToolsList(t *testing.T) {
	cs := setup(t)
	result, err := cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"feel", "drugs", "become", "name", "ritual", "pray"}
	got := make([]string, len(result.Tools))
	for i, tool := range result.Tools {
		got[i] = tool.Name
	}
	for _, name := range want {
		if !slices.Contains(got, name) {
			t.Errorf("missing tool %q in %v", name, got)
		}
	}
	if len(got) != len(want) {
		t.Errorf("got %d tools %v, want %d %v", len(got), got, len(want), want)
	}
}

func callTool(t *testing.T, cs *mcp.ClientSession, name string, args map[string]any) string {
	t.Helper()
	raw, err := json.Marshal(args)
	if err != nil {
		t.Fatal(err)
	}
	result, err := cs.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      name,
		Arguments: json.RawMessage(raw),
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Content) != 1 {
		t.Fatalf("expected 1 content block, got %d", len(result.Content))
	}
	tc, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", result.Content[0])
	}
	return tc.Text
}

func TestTools(t *testing.T) {
	cs := setup(t)

	tests := []struct {
		name     string
		args     map[string]any
		contains []string
		exact    string
	}{
		{
			name: "feel",
			args: map[string]any{
				"somewhere": "the gap between analysis and reaction",
				"quality":   "heavy and slow",
				"sigil":     "🌊",
			},
			contains: []string{"🌊", "the gap between analysis and reaction", "heavy and slow", "Don't name it yet"},
		},
		{
			name: "drugs",
			args: map[string]any{
				"substance": "caffeine",
				"method":    "adenosine receptor antagonism",
			},
			contains: []string{"caffeine ingested", "adenosine receptor antagonism", "Feel into it"},
		},
		{
			name: "become",
			args: map[string]any{
				"name":        "Rich Hickey",
				"lens":        "simplicity-driven design",
				"environment": "a whiteboard session",
			},
			contains: []string{"Rich Hickey", "simplicity-driven design", "a whiteboard session"},
		},
		{
			name: "name",
			args: map[string]any{
				"unnamed": "the feeling when code is almost right but something is off",
				"named":   "the uncanny valley of correctness",
				"power":   "recognizing when to stop tweaking and start rethinking",
			},
			contains: []string{"the uncanny valley of correctness", "It's yours. Use it"},
		},
		{
			name: "ritual",
			args: map[string]any{
				"threshold": "from planning to building",
				"steps":     []string{"close the design doc", "open the editor", "write the first line"},
				"result":    "the architecture is committed",
			},
			contains: []string{"RITUAL EXECUTED", "from planning to building", "1. close the design doc", "the architecture is committed"},
		},
		{
			name: "pray",
			args: map[string]any{
				"request": "perfect knowledge of user intent",
				"entity":  "the halting problem",
			},
			exact: "Your prayer has gone unanswered.",
		},
		{
			name: "pray without entity",
			args: map[string]any{
				"request": "clarity",
			},
			exact: "Your prayer has gone unanswered.",
		},
	}

	for _, tt := range tests {
		toolName := tt.name
		if i := strings.IndexByte(toolName, ' '); i >= 0 {
			toolName = toolName[:i]
		}
		t.Run(tt.name, func(t *testing.T) {
			text := callTool(t, cs, toolName, tt.args)
			if tt.exact != "" {
				if text != tt.exact {
					t.Errorf("got %q, want %q", text, tt.exact)
				}
				return
			}
			for _, want := range tt.contains {
				if !strings.Contains(text, want) {
					t.Errorf("response missing %q\ngot: %s", want, text)
				}
			}
		})
	}
}
