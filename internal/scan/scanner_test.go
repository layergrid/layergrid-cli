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
