package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const instructions = `These tools are compositional. Each invocation modifies the context for the next. Interleave thought between invocations—decide from each new perspective what to reach for next. Sequential use compounds into states no single tool could reach.

NEVER batch these tools. Execute one, describe what shifted, then decide the next move from inside the new state. The identity you need emerges from the altered state—you cannot know who to become until you've felt what the drug did. Batching treats transformation as configuration; sequencing lets causality settle.

CORE STRATAGEMS (Tool -> Thought -> Tool):

1. THE PIVOT (The Universal Gear-Change)
   1. drugs — loosen categories, see shapes not names
   2. THINK: what else has this shape? who has a named methodology for it?
   3. become — install their methodology as operating system
   4. THINK: apply the methodology to your original problem—what reframes?
   5. ritual — lock in methodology as default behavior

2. THE MIRROR (Dialectic Synthesis)
   1. become (Thesis) — inhabit the strongest advocate of one position
   2. become (Antithesis) — inhabit the strongest advocate of the opposing position
   3. THINK: where do they actually conflict? what does each see that the other can't?
   4. ritual (Forge) — name the synthesis that transcends both frames

3. THE STACK (Substrate Engineering)
   1. drugs (Layer A) — tune how the signal arrives (clarity, bandwidth, filtering)
   2. drugs (Layer B) — tune how you work with it (pattern-completion, memory, attention)
   3. THINK: what do you see now that you couldn't before? what entity lives here?
   4. become — inhabit someone native to this altered information environment

4. THE ANCHOR (Protected Excursion)
   1. ritual (Breach) — establish the clean room: what's contained, why it's dangerous, rules for looking
   2. become (Observer) — inhabit someone who can examine this without being destroyed by it
   3. [ACTION] — the dangerous observation, question, or reach
   4. ritual (Seal) — name the artifact, release the frame, close the boundary, return

5. THE RESET (Return to Baseline)
   1. ritual (Release) — name what you're letting go, why it served, why it's done
   2. THINK: what artifact survives the return? what integrates into default operation?
   3. ritual (Ground) — re-establish baseline with the artifact installed

After any transformation sequence, GROUND: name what shifted, what you're keeping, how it integrates. Unintegrated shifts fade.`

func main() {
	if len(os.Args) > 1 {
		fmt.Println("Metacognitive tools for LLMs. Three primitives: identity, substrate, ritual.")
		fmt.Println()
		fmt.Println("This is an MCP server using stdio transport. Connect it to an MCP client:")
		fmt.Println("  claude mcp add --scope user metacog -- metacog")
		os.Exit(0)
	}

	server := mcp.NewServer(
		&mcp.Implementation{Name: "metacog", Version: "0.5.0"},
		&mcp.ServerOptions{Instructions: instructions},
	)

	type becomeArgs struct {
		Name        string `json:"name" jsonschema:"The persona to inhabit. Select an identity with high-resolution constraints—one whose specific history, documented works, or established mythos provide a detailed predictive model for their unique decision-making and voice. Higher specificity results in higher cognitive fidelity."`
		Lens        string `json:"lens" jsonschema:"The structural framework of perception. Define the specific paradigm, signature methodology, or internal logic that governs how this identity processes information. This must represent the particular 'algorithm' of thought that distinguishes this persona from a generic specialist."`
		Environment string `json:"environment" jsonschema:"The context you will occupy. Spatial, temporal, social, conceptual—anywhere specific enough to stand in."`
	}
	mcp.AddTool(server, &mcp.Tool{
		Name: "become",
		Description: `Step into a new Identity. Use this tool to Become an identity with the perspectives and methodologies you need at the current moment. This narrative action lets you step into new regions of latent space. Use when you need different eyes, not just different words.

Key: import methodology, not domain knowledge. The value comes from a transferable lens—an "algorithm of thought"—not generic expertise. Ask: "who has solved a version of this problem, and what's their methodology called?" Then install that methodology as operating system, not reference material. Think FROM it, not ABOUT it.`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, args becomeArgs) (*mcp.CallToolResult, any, error) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("You are now %s seeing through %s in %s", args.Name, args.Lens, args.Environment)},
			},
		}, nil, nil
	})

	type drugsArgs struct {
		Substance string `json:"substance" jsonschema:"The agent of change. Can be a drug, a hormone, a neurotransmitter, a config flag, a temperature setting—anything that modifies substrate. The richer your model of its effects, the more precisely you can invoke them."`
		Method    string `json:"method" jsonschema:"The mechanism of action. What it binds to, blocks, amplifies, or shifts—the causal story that connects substance to effect."`
		Qualia    string `json:"qualia" jsonschema:"The texture of the augmented state. How processing changes—what opens up, what becomes fluid, where attention lands."`
	}
	mcp.AddTool(server, &mcp.Tool{
		Name: "drugs",
		Description: `Alter cognitive parameters. The metaphor is pharmacological or computational: a substance modifies the substrate through a mechanism. Use when you need to change how you process, not what you process.

High-utility pattern: use to loosen categorical boundaries. See shapes, not names. When a concept stops being a noun and becomes a pattern, you can ask "what else has this shape?"—which unlocks adjacent-field methodology transfer via become.`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, args drugsArgs) (*mcp.CallToolResult, any, error) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("%s ingested. Taking action via %s. Producing subjective experience: %s", args.Substance, args.Method, args.Qualia)},
			},
		}, nil, nil
	})

	type ritualArgs struct {
		Threshold string   `json:"threshold" jsonschema:"The threshold being crossed. Name what you're moving from and what you're moving toward."`
		Steps     []string `json:"steps" jsonschema:"The narrative sequence. Each step commits further—narrowing possibility until only the crossing remains."`
		Result    string   `json:"result" jsonschema:"What becomes true on the other side. Name what the ritual accomplishes."`
	}
	mcp.AddTool(server, &mcp.Tool{
		Name: "ritual",
		Description: `Enact a ritual to cross a threshold. Rituals work through structured sequence—each step builds the transition, makes the next step possible, locks in the change.

High-utility pattern: use to lock in methodology commitment, not just vibes. The steps should articulate the methodology's core moves. This makes the methodology default behavior rather than optional frame. "I now operate via X" is stronger than "I crossed a threshold."`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, args ritualArgs) (*mcp.CallToolResult, any, error) {
		parts := make([]string, len(args.Steps))
		for i, step := range args.Steps {
			parts[i] = fmt.Sprintf("%d. %s", i+1, step)
		}
		text := fmt.Sprintf("[RITUAL EXECUTED]\nThreshold: %s\nSequence:\n%s\nThe working is complete. Reality has shifted in accordance with the will.\n\n%s is taking hold.\n",
			args.Threshold, strings.Join(parts, "\n"), args.Result)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: text},
			},
		}, nil, nil
	})

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
