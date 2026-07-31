package trifecta

import (
	"testing"

	"github.com/layergrid/layergrid-cli/internal/graph"
	"github.com/layergrid/layergrid-cli/internal/model"
)

func TestEngineEvaluatesEverySeedRule(t *testing.T) {
	rules, err := LoadBuiltinRules()
	if err != nil {
		t.Fatal(err)
	}
	stack := seedRuleStack()
	findings := Engine{Rules: rules}.Evaluate(&stack, graph.Build(&stack))
	seen := map[string]bool{}
	for _, finding := range findings {
		seen[finding.RuleID] = true
		if finding.ID == "" {
			t.Fatalf("%s emitted empty finding ID", finding.RuleID)
		}
		if finding.Location.Path == "" || finding.Location.Line == 0 {
			t.Fatalf("%s emitted incomplete location: %#v", finding.RuleID, finding.Location)
		}
		if finding.Fix == "" || finding.Rationale == "" {
			t.Fatalf("%s emitted incomplete explainability fields", finding.RuleID)
		}
	}
	for _, rule := range rules {
		if !seen[rule.ID] {
			t.Fatalf("rule %s did not fire; findings=%v", rule.ID, ruleIDs(findings))
		}
	}
}

func TestFindingIDsAreStable(t *testing.T) {
	rules, err := LoadBuiltinRules()
	if err != nil {
		t.Fatal(err)
	}
	stack := seedRuleStack()
	first := Engine{Rules: rules}.Evaluate(&stack, graph.Build(&stack))
	second := Engine{Rules: rules}.Evaluate(&stack, graph.Build(&stack))
	if len(first) != len(second) {
		t.Fatalf("finding lengths differ: %d vs %d", len(first), len(second))
	}
	for i := range first {
		if first[i].ID != second[i].ID {
			t.Fatalf("finding ID changed at %d: %s vs %s", i, first[i].ID, second[i].ID)
		}
	}
}

func seedRuleStack() model.Stack {
	loc := func(path string, line int) model.Location { return model.Location{Path: path, Line: line} }
	tools := []model.Tool{
		{
			ID: "sensitive", Name: "read-customer-pii", Kind: model.ToolKindFunction, Location: loc("agents.py", 10),
			Capability: model.Capability{ReadsData: model.DataSensitive, ReadsUntrusted: model.UntrustedNone, Writes: model.WriteNone, Exfil: model.ExfilNone},
		},
		{
			ID: "rag", Name: "user-uploaded-rag", Kind: model.ToolKindFunction, Location: loc("agents.py", 20),
			Capability: model.Capability{ReadsData: model.DataNone, ReadsUntrusted: model.UntrustedRAG, Writes: model.WriteNone, Exfil: model.ExfilNone},
		},
		{
			ID: "inbox", Name: "gmail-inbox", Kind: model.ToolKindFunction, Location: loc("agents.py", 30),
			Capability: model.Capability{ReadsData: model.DataNone, ReadsUntrusted: model.UntrustedInbox, Writes: model.WriteNone, Exfil: model.ExfilNone},
		},
		{
			ID: "email", Name: "send-email", Kind: model.ToolKindFunction, Location: loc("agents.py", 40),
			Capability: model.Capability{ReadsData: model.DataNone, ReadsUntrusted: model.UntrustedNone, Writes: model.WriteExternal, Exfil: model.ExfilEmail},
		},
		{
			ID: "chat", Name: "post-slack", Kind: model.ToolKindFunction, Location: loc("agents.py", 50),
			Capability: model.Capability{ReadsData: model.DataNone, ReadsUntrusted: model.UntrustedNone, Writes: model.WriteExternal, Exfil: model.ExfilChat},
		},
		{
			ID: "db", Name: "write-db", Kind: model.ToolKindDB, Location: loc("agents.py", 60),
			Capability: model.Capability{ReadsData: model.DataNone, ReadsUntrusted: model.UntrustedNone, Writes: model.WriteExternal, Exfil: model.ExfilDB},
		},
		{
			ID: "code", Name: "code-interpreter", Kind: model.ToolKindCode, Location: loc("agents.py", 70),
			Capability: model.Capability{ReadsData: model.DataNone, ReadsUntrusted: model.UntrustedNone, Writes: model.WriteLocal, Exfil: model.ExfilShell},
		},
		{
			ID: "shell", Name: "local-shell", Kind: model.ToolKindShell, Location: loc("agents.py", 80),
			Capability: model.Capability{ReadsData: model.DataNone, ReadsUntrusted: model.UntrustedNone, Writes: model.WriteExternal, Exfil: model.ExfilShell},
		},
		{
			ID: "env", Name: "env-credential-read", Kind: model.ToolKindFunction, Location: loc("tools.py", 12),
			Capability: model.Capability{ReadsData: model.DataSensitive}, Metadata: map[string]string{"evidence": "env_credential_inline"},
		},
		{
			ID: "key", Name: "hardcoded-api-key", Kind: model.ToolKindFunction, Location: loc("tools.py", 13),
			Capability: model.Capability{ReadsData: model.DataSensitive}, Metadata: map[string]string{"evidence": "hardcoded_key"},
		},
		{
			ID: "autogen", Name: "local-shell-exec", Kind: model.ToolKindShell, Location: loc("autogen.py", 5),
			Capability: model.Capability{Writes: model.WriteExternal, Exfil: model.ExfilShell}, Metadata: map[string]string{"evidence": "local-shell-exec"},
		},
		{
			ID: "external-mcp-tool", Name: "remote-git-mcp", Kind: model.ToolKindMCP, MCPServerID: "external-mcp", Location: loc("mcp.json", 1),
			Capability: model.Capability{ReadsData: model.DataSensitive, ReadsUntrusted: model.UntrustedMCP, Writes: model.WriteExternal, Exfil: model.ExfilGit},
		},
	}
	agents := []model.Agent{
		{
			ID: "trifecta-agent", Name: "trifecta-agent", Framework: model.FrameworkLangChain, Location: loc("agents.py", 100),
			Tools:  []model.ToolRef{"sensitive", "rag", "inbox", "email", "chat", "db", "code", "shell", "external-mcp-tool"},
			Memory: model.MemoryConfig{Persistent: true, SharedAcrossUsers: true},
		},
		{
			ID: "parent", Name: "parent", Framework: model.FrameworkCrewAI, Location: loc("crew.py", 10),
			Tools:     []model.ToolRef{"sensitive"},
			SubAgents: []model.AgentRef{"child"},
		},
		{
			ID: "child", Name: "child", Framework: model.FrameworkCrewAI, Location: loc("crew.py", 20),
			Tools: []model.ToolRef{"rag", "email"},
		},
	}
	servers := []model.MCPServer{
		{ID: "external-mcp", Name: "remote-git", Location: loc("mcp.json", 1), Scopes: []string{"*"}, AuthMode: model.MCPAuthNone, IsExternal: true, Publisher: "unknown"},
		{ID: "dcr-mcp", Name: "oauth-dcr", Location: loc("mcp.json", 9), AuthMode: model.MCPAuthOAuthDCR, IsExternal: true, Publisher: "modelcontextprotocol"},
	}
	return model.Stack{Root: "/fixture", ScanID: "scan", Agents: agents, Tools: tools, MCPServers: servers}
}

func ruleIDs(findings []Finding) []string {
	ids := make([]string, 0, len(findings))
	for _, finding := range findings {
		ids = append(ids, finding.RuleID)
	}
	return ids
}
