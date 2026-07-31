package trifecta

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/layergrid/layergrid-cli/internal/graph"
	"github.com/layergrid/layergrid-cli/internal/model"
)

type Engine struct {
	Rules []Rule
}

func (e Engine) Evaluate(s *model.Stack, g *graph.Graph) []Finding {
	toolsByID := map[string]model.Tool{}
	for _, tool := range s.Tools {
		toolsByID[tool.ID] = tool
	}
	var findings []Finding
	for _, rule := range e.Rules {
		findings = append(findings, evaluateRule(rule, s, toolsByID, g)...)
	}
	sort.Slice(findings, func(i, j int) bool {
		if findings[i].Severity != findings[j].Severity {
			return severityRank(findings[i].Severity) < severityRank(findings[j].Severity)
		}
		if findings[i].RuleID != findings[j].RuleID {
			return findings[i].RuleID < findings[j].RuleID
		}
		return findings[i].ID < findings[j].ID
	})
	return findings
}

func evaluateRule(rule Rule, s *model.Stack, toolsByID map[string]model.Tool, _ *graph.Graph) []Finding {
	var out []Finding
	for _, agent := range s.Agents {
		caps := capabilities(agent, toolsByID)
		switch rule.ID {
		case "LG-LETHAL-TRIFECTA-01":
			if hasSensitive(caps) && hasUntrusted(caps) && hasExfil(caps) {
				out = append(out, finding(rule, "agent", agent.ID, agent.Name, agent.Location, pathForAgent(agent, toolsByID), "agent combines sensitive data access, untrusted input, and external communication"))
			}
		case "LG-TOOL-EXFIL-EMAIL-01":
			if hasUntrustedInbox(caps) && hasExfilKind(caps, model.ExfilEmail) {
				out = append(out, finding(rule, "agent", agent.ID, agent.Name, agent.Location, pathForAgent(agent, toolsByID), "agent can read inbox content and send email"))
			}
		case "LG-TOOL-EXFIL-CHAT-01":
			if hasSensitive(caps) && hasExfilKind(caps, model.ExfilChat) {
				out = append(out, finding(rule, "agent", agent.ID, agent.Name, agent.Location, pathForAgent(agent, toolsByID), "agent can read sensitive data and post to chat"))
			}
		case "LG-TOOL-EXFIL-DBWRITE-01":
			if hasSensitive(caps) && hasExfilKind(caps, model.ExfilDB) {
				out = append(out, finding(rule, "agent", agent.ID, agent.Name, agent.Location, pathForAgent(agent, toolsByID), "agent can read sensitive data and write to an external database"))
			}
		case "LG-TOOL-CODE-EXEC-01":
			if hasToolKind(agent, toolsByID, model.ToolKindCode) {
				out = append(out, finding(rule, "agent", agent.ID, agent.Name, agent.Location, pathForAgent(agent, toolsByID), "agent has a code execution tool"))
			}
		case "LG-TOOL-SHELL-EXEC-01":
			if hasToolKind(agent, toolsByID, model.ToolKindShell) {
				out = append(out, finding(rule, "agent", agent.ID, agent.Name, agent.Location, pathForAgent(agent, toolsByID), "agent has shell execution capability"))
			}
		case "LG-MEMORY-UNBOUNDED-01":
			if agent.Memory.Persistent && agent.Memory.RetentionPolicy == "" {
				out = append(out, finding(rule, "agent", agent.ID, agent.Name, agent.Location, nil, "agent has persistent memory without a retention policy"))
			}
		case "LG-MEMORY-CROSS-USER-01":
			if agent.Memory.SharedAcrossUsers {
				out = append(out, finding(rule, "agent", agent.ID, agent.Name, agent.Location, nil, "agent memory appears to be shared across users"))
			}
		case "LG-LETHAL-TRIFECTA-02":
			if agent.Memory.SharedAcrossUsers && stackHasSensitive(s) && stackHasUntrusted(s) && stackHasExfil(s) {
				out = append(out, finding(rule, "agent", agent.ID, agent.Name, agent.Location, nil, "shared-memory agent group collectively combines lethal-trifecta capabilities"))
			}
		case "LG-LETHAL-TRIFECTA-03":
			if len(agent.SubAgents) > 0 && hasTrifectaAcrossCascade(agent, toolsByID, s.Agents) {
				out = append(out, finding(rule, "agent", agent.ID, agent.Name, agent.Location, nil, "agent handoff chain collectively combines lethal-trifecta capabilities"))
			}
		case "LG-MCP-EXTERNAL-WRITE-01":
			if hasSensitive(caps) && hasExternalMCPWrite(agent, toolsByID, s.MCPServers) {
				out = append(out, finding(rule, "agent", agent.ID, agent.Name, agent.Location, pathForAgent(agent, toolsByID), "agent combines sensitive data access with an external MCP write channel"))
			}
		case "LG-RAG-UNTRUSTED-01":
			if hasRAG(caps) && len(agent.Guardrails) == 0 {
				out = append(out, finding(rule, "agent", agent.ID, agent.Name, agent.Location, pathForAgent(agent, toolsByID), "agent consumes RAG or user-uploaded content without a detected sanitization guardrail"))
			}
		case "LG-AGENT-NO-GUARDRAIL-01":
			if len(agent.Guardrails) == 0 && hasExfil(caps) {
				out = append(out, finding(rule, "agent", agent.ID, agent.Name, agent.Location, nil, "agent has external tools but no detected guardrail middleware"))
			}
		}
	}
	for _, tool := range s.Tools {
		switch rule.ID {
		case "LG-CREDENTIAL-ENV-IN-CONTEXT-01":
			if tool.Metadata["evidence"] == "env_credential_inline" {
				out = append(out, finding(rule, "tool", tool.ID, tool.Name, tool.Location, []PathNode{{Kind: "tool", ID: tool.ID, Name: tool.Name}}, "tool reads environment credentials inline"))
			}
		case "LG-CREDENTIAL-KEY-HARDCODED-01":
			if tool.Metadata["evidence"] == "hardcoded_key" {
				out = append(out, finding(rule, "tool", tool.ID, tool.Name, tool.Location, []PathNode{{Kind: "tool", ID: tool.ID, Name: tool.Name}}, "credential-like literal detected in agent code"))
			}
		case "LG-AUTOGEN-LOCAL-EXEC-01":
			if tool.Metadata["evidence"] == "local-shell-exec" {
				out = append(out, finding(rule, "tool", tool.ID, tool.Name, tool.Location, []PathNode{{Kind: "tool", ID: tool.ID, Name: tool.Name}}, "AutoGen LocalCommandLineCodeExecutor detected"))
			}
		}
	}
	for _, server := range s.MCPServers {
		switch rule.ID {
		case "LG-MCP-OVERSCOPE-01":
			if hasScope(server.Scopes, "*") {
				out = append(out, finding(rule, "mcp", server.ID, server.Name, server.Location, nil, "MCP server is registered with wildcard scope"))
			}
		case "LG-MCP-NOAUTH-01":
			if server.IsExternal && server.AuthMode == model.MCPAuthNone {
				out = append(out, finding(rule, "mcp", server.ID, server.Name, server.Location, nil, "external MCP server has no detected auth mode"))
			}
		case "LG-MCP-PUBLISHER-UNKNOWN-01":
			if server.Publisher == "" || server.Publisher == "unknown" {
				out = append(out, finding(rule, "mcp", server.ID, server.Name, server.Location, nil, "MCP server publisher could not be verified from config"))
			}
		case "LG-MCP-DCR-01":
			if server.AuthMode == model.MCPAuthOAuthDCR {
				out = append(out, finding(rule, "mcp", server.ID, server.Name, server.Location, nil, "MCP server uses OAuth Dynamic Client Registration"))
			}
		}
	}
	return out
}

func stackHasSensitive(s *model.Stack) bool {
	for _, tool := range s.Tools {
		if tool.Capability.ReadsData == model.DataSensitive || tool.Capability.ReadsData == model.DataRestricted {
			return true
		}
	}
	return false
}

func stackHasUntrusted(s *model.Stack) bool {
	for _, tool := range s.Tools {
		if tool.Capability.ReadsUntrusted != "" && tool.Capability.ReadsUntrusted != model.UntrustedNone {
			return true
		}
	}
	return false
}

func stackHasExfil(s *model.Stack) bool {
	for _, tool := range s.Tools {
		if tool.Capability.Exfil != "" && tool.Capability.Exfil != model.ExfilNone {
			return true
		}
	}
	return false
}

func hasTrifectaAcrossCascade(agent model.Agent, toolsByID map[string]model.Tool, agents []model.Agent) bool {
	allCaps := capabilities(agent, toolsByID)
	children := map[string]bool{}
	for _, ref := range agent.SubAgents {
		children[string(ref)] = true
	}
	for _, candidate := range agents {
		if children[candidate.ID] {
			allCaps = append(allCaps, capabilities(candidate, toolsByID)...)
		}
	}
	return hasSensitive(allCaps) && hasUntrusted(allCaps) && hasExfil(allCaps)
}

func hasExternalMCPWrite(agent model.Agent, toolsByID map[string]model.Tool, servers []model.MCPServer) bool {
	externalServers := map[string]bool{}
	for _, server := range servers {
		if server.IsExternal {
			externalServers[server.ID] = true
		}
	}
	for _, ref := range agent.Tools {
		tool, ok := toolsByID[string(ref)]
		if !ok || tool.Kind != model.ToolKindMCP {
			continue
		}
		if externalServers[tool.MCPServerID] && tool.Capability.Writes == model.WriteExternal {
			return true
		}
	}
	return false
}

func capabilities(agent model.Agent, toolsByID map[string]model.Tool) []model.Capability {
	out := make([]model.Capability, 0, len(agent.Tools))
	for _, ref := range agent.Tools {
		if tool, ok := toolsByID[string(ref)]; ok {
			out = append(out, tool.Capability)
		}
	}
	return out
}

func pathForAgent(agent model.Agent, toolsByID map[string]model.Tool) []PathNode {
	out := []PathNode{{Kind: "agent", ID: agent.ID, Name: agent.Name}}
	for _, ref := range agent.Tools {
		if tool, ok := toolsByID[string(ref)]; ok {
			out = append(out, PathNode{Kind: "tool", ID: tool.ID, Name: tool.Name})
		}
	}
	return out
}

func finding(rule Rule, kind, id, name string, loc model.Location, path []PathNode, rationale string) Finding {
	f := Finding{
		RuleID: rule.ID, RuleName: rule.Name, Severity: rule.Severity,
		Category: rule.Category,
		Subject:  Subject{Kind: kind, ID: id, Name: name}, Path: path,
		Location: loc, Fix: rule.Fix, References: rule.References,
		ScoreImpact: rule.ScoreImpact, Rationale: rationale,
		Confidence: confidenceFor(loc),
	}
	if f.Confidence == "low" {
		f.ScoreImpact = 0
	}
	f.ID = stableFindingID(rule.ID, id, loc)
	return f
}

func confidenceFor(loc model.Location) string {
	path := strings.ToLower(loc.Path)
	lowSignals := []string{"/test/", "/tests/", "test_", "_test.", "/examples/", "/example/", "/docs/", "/doc/", "cookbook", "/notebooks/", ".ipynb"}
	for _, signal := range lowSignals {
		if strings.Contains(path, signal) {
			return "low"
		}
	}
	return "high"
}

func stableFindingID(ruleID, subjectID string, loc model.Location) string {
	h := sha256.New()
	_, _ = fmt.Fprintf(h, "%s|%s|%s|%d", ruleID, subjectID, loc.Path, loc.Line)
	return hex.EncodeToString(h.Sum(nil))[:24]
}

func hasSensitive(caps []model.Capability) bool {
	for _, cap := range caps {
		if cap.ReadsData == model.DataSensitive || cap.ReadsData == model.DataRestricted {
			return true
		}
	}
	return false
}

func hasUntrusted(caps []model.Capability) bool {
	for _, cap := range caps {
		if cap.ReadsUntrusted != "" && cap.ReadsUntrusted != model.UntrustedNone {
			return true
		}
	}
	return false
}

func hasUntrustedInbox(caps []model.Capability) bool {
	for _, cap := range caps {
		if cap.ReadsUntrusted == model.UntrustedInbox {
			return true
		}
	}
	return false
}

func hasRAG(caps []model.Capability) bool {
	for _, cap := range caps {
		if cap.ReadsUntrusted == model.UntrustedRAG {
			return true
		}
	}
	return false
}

func hasExfil(caps []model.Capability) bool {
	for _, cap := range caps {
		if cap.Exfil != "" && cap.Exfil != model.ExfilNone {
			return true
		}
	}
	return false
}

func hasExfilKind(caps []model.Capability, exfil model.ExfilChannel) bool {
	for _, cap := range caps {
		if cap.Exfil == exfil {
			return true
		}
	}
	return false
}

func hasToolKind(agent model.Agent, toolsByID map[string]model.Tool, kind model.ToolKind) bool {
	for _, ref := range agent.Tools {
		if tool, ok := toolsByID[string(ref)]; ok && tool.Kind == kind {
			return true
		}
	}
	return false
}

func hasScope(scopes []string, target string) bool {
	for _, scope := range scopes {
		if scope == target {
			return true
		}
	}
	return false
}

func severityRank(s Severity) int {
	switch s {
	case SeverityCritical:
		return 0
	case SeverityHigh:
		return 1
	case SeverityMedium:
		return 2
	default:
		return 3
	}
}
