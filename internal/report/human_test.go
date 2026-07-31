package report

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/layergrid/layergrid-cli/internal/model"
	"github.com/layergrid/layergrid-cli/internal/scan"
	"github.com/layergrid/layergrid-cli/internal/trifecta"
)

func TestHumanOutputGolden(t *testing.T) {
	result := scan.Result{
		Stack: model.Stack{
			Root: "/Users/nikhil/dev/my-app",
			Agents: []model.Agent{
				{ID: "support-agent", Name: "support-agent", Framework: model.FrameworkLangChain},
				{ID: "intake-agent", Name: "intake-agent", Framework: model.FrameworkCrewAI},
			},
			Tools: []model.Tool{
				{Kind: model.ToolKindFunction}, {Kind: model.ToolKindMCP}, {Kind: model.ToolKindShell},
			},
			MCPServers:  []model.MCPServer{{IsExternal: true, Publisher: "unknown"}},
			Datasources: []model.Datasource{{ID: "prod-db"}},
		},
		Findings: []trifecta.Finding{
			{
				RuleID: "LG-LETHAL-TRIFECTA-01", Severity: trifecta.SeverityCritical, ScoreImpact: -30,
				Path:     []trifecta.PathNode{{Kind: "agent", Name: "support-agent"}, {Kind: "mcp", Name: "Slack MCP"}, {Kind: "data", Name: "Secrets store"}},
				Location: model.Location{Path: "internal/mcp/slack.py", Line: 41},
				Fix:      "scope Slack MCP to chat:read only",
			},
			{
				RuleID: "LG-MCP-DCR-01", Severity: trifecta.SeverityMedium, ScoreImpact: -5,
				Subject:   trifecta.Subject{Name: "ops-agent"},
				Rationale: "ops-agent uses OAuth DCR auth mode",
				Location:  model.Location{Path: "mcp.json", Line: 2},
				Fix:       "pin to static OAuth registration",
			},
		},
		Score: trifecta.Score{
			Value: 65, Grade: "C",
			Counts: map[trifecta.Severity]int{trifecta.SeverityCritical: 1, trifecta.SeverityHigh: 0, trifecta.SeverityMedium: 1, trifecta.SeverityLow: 0},
		},
		Duration: 4200 * time.Millisecond,
	}
	got, err := (Human{}).Format(result)
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(filepath.Join("..", "..", "testdata", "golden", "human-output.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Fatalf("human output mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestHumanNoColorUsesASCII(t *testing.T) {
	result := scan.Result{
		Stack: model.Stack{Root: "/repo", Agents: []model.Agent{{Framework: model.FrameworkLangChain}}, Tools: []model.Tool{}, MCPServers: []model.MCPServer{}, Datasources: []model.Datasource{}},
		Score: trifecta.Score{Value: 100, Grade: "A", Counts: map[trifecta.Severity]int{trifecta.SeverityCritical: 0, trifecta.SeverityHigh: 0, trifecta.SeverityMedium: 0, trifecta.SeverityLow: 0}},
	}
	got, err := (Human{NoColor: true}).Format(result)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range string(got) {
		if r > 127 {
			t.Fatalf("no-color output contains non-ASCII rune %q in:\n%s", r, got)
		}
	}
}
