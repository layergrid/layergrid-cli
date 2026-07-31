package trifecta

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/layergrid/layergrid/internal/graph"
	"github.com/layergrid/layergrid/internal/model"
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
		case "LG-AGENT-NO-GUARDRAIL-01":
			if len(agent.Guardrails) == 0 && hasExfil(caps) {
				out = append(out, finding(rule, "agent", agent.ID, agent.Name, agent.Location, nil, "agent has external tools but no detected guardrail middleware"))
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
		Subject: Subject{Kind: kind, ID: id, Name: name}, Path: path,
		Location: loc, Fix: rule.Fix, References: rule.References,
		ScoreImpact: rule.ScoreImpact, Rationale: rationale,
	}
	f.ID = stableFindingID(rule.ID, id, loc)
	return f
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
