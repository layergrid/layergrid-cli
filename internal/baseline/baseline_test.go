package baseline

import (
	"testing"
	"time"

	"github.com/layergrid/layergrid-cli/internal/model"
)

func TestFromStackIsDeterministic(t *testing.T) {
	stack := model.Stack{
		Agents:     []model.Agent{{ID: "agent", Name: "agent", Tools: []model.ToolRef{"tool"}}},
		Tools:      []model.Tool{{ID: "tool", Name: "send", Description: "send message", Scope: []string{"write", "read"}}},
		MCPServers: []model.MCPServer{{ID: "mcp", Name: "remote", Transport: "http", Endpoint: "https://example.com", Scopes: []string{"write", "read"}}},
	}
	captured := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	first := FromStack(stack, captured)
	second := FromStack(stack, captured)
	if first.Tools[0].Descriptor != second.Tools[0].Descriptor {
		t.Fatalf("tool descriptor changed: %s vs %s", first.Tools[0].Descriptor, second.Tools[0].Descriptor)
	}
	if first.MCPServers[0].Descriptor != second.MCPServers[0].Descriptor {
		t.Fatalf("mcp descriptor changed: %s vs %s", first.MCPServers[0].Descriptor, second.MCPServers[0].Descriptor)
	}
}

func TestCompareDetectsScopeWideningAndToolAdded(t *testing.T) {
	before := Baseline{
		CapturedAt: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
		Tools:      []ToolRecord{{ID: "tool", Name: "tool", Scopes: []string{"read"}, Descriptor: "sha256:a"}},
		Agents:     []AgentRecord{{ID: "agent", Name: "agent", Tools: []string{"tool"}}},
	}
	after := Baseline{
		Tools:  []ToolRecord{{ID: "tool", Name: "tool", Scopes: []string{"read", "write"}, Descriptor: "sha256:a"}, {ID: "new", Name: "new"}},
		Agents: []AgentRecord{{ID: "agent", Name: "agent", Tools: []string{"new", "tool"}}},
	}
	result := Compare(before, after)
	if result.Summary.ScopeWidening != 1 || result.Summary.ToolAdded != 1 {
		t.Fatalf("summary = %#v, want one scope widening and one tool added", result.Summary)
	}
	if !ShouldFail(result, "scope-widening") || !ShouldFail(result, "tool-added") || !ShouldFail(result, "any") {
		t.Fatal("expected compare result to fail selected thresholds")
	}
}
