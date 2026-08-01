package report

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/layergrid/layergrid-cli/internal/model"
	"github.com/layergrid/layergrid-cli/internal/scan"
	"github.com/layergrid/layergrid-cli/internal/trifecta"
)

func TestHumanOutputGolden(t *testing.T) {
	got := renderHuman(t, sampleResult(), 100)
	assertGolden(t, "human-output.txt", got)
}

func TestHumanEmptyState(t *testing.T) {
	result := sampleResult()
	result.Findings = nil
	result.Score = trifecta.Score{Value: 100, Grade: "A", Counts: severityCounts(0, 0, 0, 0)}
	got := renderHuman(t, result, 100)
	assertGolden(t, "human-empty.txt", got)
	if strings.Contains(got, "Attack Paths") {
		t.Fatalf("empty output should not include attack paths:\n%s", got)
	}
}

func TestHumanSoloAgentAttackPath(t *testing.T) {
	result := sampleResult()
	result.Findings = []trifecta.Finding{{
		RuleID: "LG-TOOL-CODE-EXEC-01", Severity: trifecta.SeverityHigh, ScoreImpact: -15,
		Subject:   trifecta.Subject{Name: "autogen-runner"},
		Location:  model.Location{Path: "agents/runner.py", Line: 22},
		Fix:       "sandbox or remove code executor",
		Rationale: "agent can execute local code",
	}}
	result.Score = trifecta.Score{Value: 85, Grade: "B", Counts: severityCounts(0, 1, 0, 0)}
	got := renderHuman(t, result, 100)
	assertGolden(t, "human-solo.txt", got)
	if !strings.Contains(got, "autogen-runner") {
		t.Fatalf("solo-agent path missing subject:\n%s", got)
	}
}

func TestHumanNarrowLayout(t *testing.T) {
	got := renderHuman(t, sampleResult(), 70)
	assertGolden(t, "human-narrow.txt", got)
	if !strings.Contains(got, "\n  Breakdown\n") {
		t.Fatalf("narrow output should stack breakdown below score:\n%s", got)
	}
}

func TestHumanNoColorStripsANSI(t *testing.T) {
	got := renderHuman(t, sampleResult(), 100)
	if strings.Contains(got, "\x1b[") {
		t.Fatalf("no-color output contains ANSI escapes:\n%s", got)
	}
	if !strings.Contains(got, "┌") || !strings.Contains(got, "╭") {
		t.Fatalf("no-color output should keep Unicode table borders:\n%s", got)
	}
}

func renderHuman(t *testing.T, result scan.Result, width int) string {
	t.Helper()
	got, err := (Human{NoColor: true, Width: width}).Format(result)
	if err != nil {
		t.Fatal(err)
	}
	return string(got)
}

func assertGolden(t *testing.T, name, got string) {
	t.Helper()
	path := filepath.Join("..", "..", "testdata", "golden", name)
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != string(want) {
		t.Fatalf("%s mismatch\n--- got ---\n%s\n--- want ---\n%s", name, got, want)
	}
}

func sampleResult() scan.Result {
	return scan.Result{
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
		Score:    trifecta.Score{Value: 65, Grade: "C", Counts: severityCounts(1, 0, 1, 0)},
		Duration: 4200 * time.Millisecond,
	}
}

func severityCounts(critical, high, medium, low int) map[trifecta.Severity]int {
	return map[trifecta.Severity]int{
		trifecta.SeverityCritical: critical,
		trifecta.SeverityHigh:     high,
		trifecta.SeverityMedium:   medium,
		trifecta.SeverityLow:      low,
	}
}
