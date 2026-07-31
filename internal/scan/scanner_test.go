package scan

import (
	"path/filepath"
	"testing"
)

func TestRunLangChainTrifectaFixture(t *testing.T) {
	result, err := Run(filepath.Join("..", "..", "testdata", "langchain-trifecta"), Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Stack.Agents) != 1 {
		t.Fatalf("agents = %d, want 1", len(result.Stack.Agents))
	}
	if len(result.Findings) == 0 {
		t.Fatal("expected at least one finding")
	}
	if result.Score.Value >= 100 {
		t.Fatalf("score = %d, want below 100", result.Score.Value)
	}
}

func TestRunMCPFixture(t *testing.T) {
	result, err := Run(filepath.Join("..", "..", "testdata", "mcp-github-lethal"), Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Stack.MCPServers) != 2 {
		t.Fatalf("mcp servers = %d, want 2", len(result.Stack.MCPServers))
	}
	var foundWildcard bool
	for _, f := range result.Findings {
		if f.RuleID == "LG-MCP-OVERSCOPE-01" {
			foundWildcard = true
		}
	}
	if !foundWildcard {
		t.Fatal("expected wildcard MCP scope finding")
	}
}

func TestRunAutoGenLocalExecFixture(t *testing.T) {
	result, err := Run(filepath.Join("..", "..", "testdata", "autogen-local-exec"), Options{})
	if err != nil {
		t.Fatal(err)
	}
	assertFinding(t, result, "LG-AUTOGEN-LOCAL-EXEC-01")
	assertFinding(t, result, "LG-TOOL-SHELL-EXEC-01")
}

func TestRunOpenAICodeInterpreterFixture(t *testing.T) {
	result, err := Run(filepath.Join("..", "..", "testdata", "openai-code-interpreter"), Options{})
	if err != nil {
		t.Fatal(err)
	}
	assertFinding(t, result, "LG-TOOL-CODE-EXEC-01")
}

func TestRunConfigCanDisableRule(t *testing.T) {
	result, err := Run(filepath.Join("..", "..", "testdata", "config-disable"), Options{})
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range result.Findings {
		if f.RuleID == "LG-AGENT-NO-GUARDRAIL-01" {
			t.Fatal("expected config to disable LG-AGENT-NO-GUARDRAIL-01")
		}
	}
	assertFinding(t, result, "LG-LETHAL-TRIFECTA-01")
}

func TestRunHonorsExcludePatterns(t *testing.T) {
	result, err := Run(filepath.Join("..", "..", "testdata", "langchain-trifecta"), Options{Exclude: []string{"agent.py"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Stack.Agents) != 0 {
		t.Fatalf("agents = %d, want 0", len(result.Stack.Agents))
	}
	if len(result.Findings) != 0 {
		t.Fatalf("findings = %d, want 0", len(result.Findings))
	}
}

func assertFinding(t *testing.T, result Result, ruleID string) {
	t.Helper()
	for _, f := range result.Findings {
		if f.RuleID == ruleID {
			return
		}
	}
	t.Fatalf("expected finding %s, got %#v", ruleID, result.Findings)
}
