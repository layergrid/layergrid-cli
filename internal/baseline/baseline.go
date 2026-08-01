package baseline

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/layergrid/layergrid-cli/internal/model"
	"github.com/layergrid/layergrid-cli/internal/version"
)

const SchemaVersion = "1.0.0"

type Baseline struct {
	SchemaVersion string        `json:"schemaVersion"`
	CapturedAt    time.Time     `json:"capturedAt"`
	ToolVersion   string        `json:"toolVersion"`
	MCPServers    []MCPRecord   `json:"mcpServers"`
	Tools         []ToolRecord  `json:"tools"`
	Agents        []AgentRecord `json:"agents"`
}

type MCPRecord struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Transport  string   `json:"transport"`
	Endpoint   string   `json:"endpoint,omitempty"`
	Scopes     []string `json:"scopes,omitempty"`
	Descriptor string   `json:"descriptor"`
}

type ToolRecord struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Scopes      []string `json:"scopes,omitempty"`
	Descriptor  string   `json:"descriptor"`
}

type AgentRecord struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Tools []string `json:"tools,omitempty"`
}

type Change struct {
	Type   string `json:"type"`
	ID     string `json:"id"`
	Name   string `json:"name,omitempty"`
	Detail string `json:"detail,omitempty"`
}

type CompareResult struct {
	BaselineCapturedAt time.Time `json:"baselineCapturedAt"`
	Changes            []Change  `json:"changes"`
	Summary            Summary   `json:"summary"`
}

type Summary struct {
	Total           int `json:"total"`
	NewMCPServers   int `json:"newMcpServers"`
	ToolAdded       int `json:"toolAdded"`
	DescriptorDrift int `json:"descriptorDrift"`
	ScopeWidening   int `json:"scopeWidening"`
}

func FromStack(stack model.Stack, capturedAt time.Time) Baseline {
	b := Baseline{
		SchemaVersion: SchemaVersion,
		CapturedAt:    capturedAt.UTC(),
		ToolVersion:   version.Version,
		MCPServers:    make([]MCPRecord, 0, len(stack.MCPServers)),
		Tools:         make([]ToolRecord, 0, len(stack.Tools)),
		Agents:        make([]AgentRecord, 0, len(stack.Agents)),
	}
	for _, server := range stack.MCPServers {
		scopes := sortedCopy(server.Scopes)
		b.MCPServers = append(b.MCPServers, MCPRecord{
			ID:         server.ID,
			Name:       server.Name,
			Transport:  server.Transport,
			Endpoint:   server.Endpoint,
			Scopes:     scopes,
			Descriptor: descriptor(server.Name, "", scopes),
		})
	}
	for _, tool := range stack.Tools {
		scopes := sortedCopy(tool.Scope)
		b.Tools = append(b.Tools, ToolRecord{
			ID:          tool.ID,
			Name:        tool.Name,
			Description: tool.Description,
			Scopes:      scopes,
			Descriptor:  descriptor(tool.Name, tool.Description, scopes),
		})
	}
	for _, agent := range stack.Agents {
		tools := make([]string, 0, len(agent.Tools))
		for _, ref := range agent.Tools {
			tools = append(tools, string(ref))
		}
		sort.Strings(tools)
		b.Agents = append(b.Agents, AgentRecord{ID: agent.ID, Name: agent.Name, Tools: tools})
	}
	sort.Slice(b.MCPServers, func(i, j int) bool { return b.MCPServers[i].ID < b.MCPServers[j].ID })
	sort.Slice(b.Tools, func(i, j int) bool { return b.Tools[i].ID < b.Tools[j].ID })
	sort.Slice(b.Agents, func(i, j int) bool { return b.Agents[i].ID < b.Agents[j].ID })
	return b
}

func Save(path string, b Baseline) error {
	if existing, err := Load(path); err == nil && !existing.CapturedAt.IsZero() {
		b.CapturedAt = existing.CapturedAt
	}
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func Load(path string) (Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Baseline{}, err
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return Baseline{}, err
	}
	return b, nil
}

func Compare(before, after Baseline) CompareResult {
	result := CompareResult{BaselineCapturedAt: before.CapturedAt, Changes: []Change{}}
	beforeMCP := mapMCP(before.MCPServers)
	beforeTools := mapTools(before.Tools)
	beforeAgents := mapAgents(before.Agents)
	for _, server := range after.MCPServers {
		old, ok := beforeMCP[server.ID]
		if !ok {
			result.add(Change{Type: "mcp_added", ID: server.ID, Name: server.Name, Detail: "new MCP server"})
			continue
		}
		if old.Descriptor != server.Descriptor {
			result.add(Change{Type: "descriptor_drift", ID: server.ID, Name: server.Name, Detail: "MCP server descriptor changed"})
		}
		if widened(old.Scopes, server.Scopes) {
			result.add(Change{Type: "scope_widened", ID: server.ID, Name: server.Name, Detail: "MCP server scopes widened"})
		}
	}
	for _, tool := range after.Tools {
		old, ok := beforeTools[tool.ID]
		if !ok {
			continue
		}
		if old.Descriptor != tool.Descriptor {
			result.add(Change{Type: "descriptor_drift", ID: tool.ID, Name: tool.Name, Detail: "tool descriptor changed"})
		}
		if widened(old.Scopes, tool.Scopes) {
			result.add(Change{Type: "scope_widened", ID: tool.ID, Name: tool.Name, Detail: "tool scopes widened"})
		}
	}
	for _, agent := range after.Agents {
		old, ok := beforeAgents[agent.ID]
		if !ok {
			continue
		}
		oldTools := stringSet(old.Tools)
		for _, toolID := range agent.Tools {
			if !oldTools[toolID] {
				result.add(Change{Type: "tool_added", ID: agent.ID, Name: agent.Name, Detail: fmt.Sprintf("new tool %s added to existing agent", toolID)})
			}
		}
	}
	return result
}

func (r *CompareResult) add(change Change) {
	r.Changes = append(r.Changes, change)
	r.Summary.Total++
	switch change.Type {
	case "mcp_added":
		r.Summary.NewMCPServers++
	case "tool_added":
		r.Summary.ToolAdded++
	case "descriptor_drift":
		r.Summary.DescriptorDrift++
	case "scope_widened":
		r.Summary.ScopeWidening++
	}
}

func ShouldFail(result CompareResult, failOn string) bool {
	switch failOn {
	case "", "never":
		return false
	case "any":
		return len(result.Changes) > 0
	case "tool-added":
		return result.Summary.ToolAdded > 0
	case "descriptor-drift":
		return result.Summary.DescriptorDrift > 0
	case "scope-widening":
		return result.Summary.ScopeWidening > 0
	default:
		return false
	}
}

func FormatHuman(result CompareResult) string {
	if len(result.Changes) == 0 {
		return "No baseline drift detected.\n"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Baseline drift detected (%d changes)\n\n", len(result.Changes))
	for _, change := range result.Changes {
		fmt.Fprintf(&b, "- %s: %s", change.Type, change.ID)
		if change.Name != "" {
			fmt.Fprintf(&b, " (%s)", change.Name)
		}
		if change.Detail != "" {
			fmt.Fprintf(&b, " - %s", change.Detail)
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func descriptor(name, description string, scopes []string) string {
	parts := []string{name, description}
	parts = append(parts, sortedCopy(scopes)...)
	h := sha256.New()
	for _, part := range parts {
		_, _ = h.Write([]byte(part))
		_, _ = h.Write([]byte{0})
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil))
}

func sortedCopy(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	out := append([]string(nil), values...)
	sort.Strings(out)
	return out
}

func widened(before, after []string) bool {
	old := stringSet(before)
	for _, value := range after {
		if !old[value] {
			return true
		}
	}
	return false
}

func stringSet(values []string) map[string]bool {
	set := map[string]bool{}
	for _, value := range values {
		set[value] = true
	}
	return set
}

func mapMCP(values []MCPRecord) map[string]MCPRecord {
	out := map[string]MCPRecord{}
	for _, value := range values {
		out[value.ID] = value
	}
	return out
}

func mapTools(values []ToolRecord) map[string]ToolRecord {
	out := map[string]ToolRecord{}
	for _, value := range values {
		out[value.ID] = value
	}
	return out
}

func mapAgents(values []AgentRecord) map[string]AgentRecord {
	out := map[string]AgentRecord{}
	for _, value := range values {
		out[value.ID] = value
	}
	return out
}
