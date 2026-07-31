package claude_sdk

import (
	"strings"

	"github.com/layergrid/layergrid-cli/internal/detectors/pyutil"
	"github.com/layergrid/layergrid-cli/internal/model"
)

type Detector struct{}

func New() Detector { return Detector{} }

func (Detector) Name() string               { return "claude_sdk" }
func (Detector) Framework() model.Framework { return model.FrameworkClaudeSDK }

func (Detector) Detect(root string, s *model.Stack) error {
	return pyutil.Walk(root, func(path string, lines []string) error {
		if !pyutil.HasAny(lines, "claude_agent_sdk", "AgentDefinition", "allowed_tools", "mcp_servers") {
			return nil
		}
		var tools []model.ToolRef
		for i, line := range lines {
			lower := strings.ToLower(line)
			if strings.Contains(lower, "allowed_tools") || strings.Contains(lower, "mcp_servers") {
				tool := model.Tool{
					ID:         model.StableID("tool", "claude-sdk", path, line),
					Name:       "claude-allowed-tools",
					Kind:       model.ToolKindMCP,
					Source:     model.ToolSource{Kind: "python", Name: path},
					Location:   model.RelativeLocation(root, path, i+1),
					Capability: inferCapability(lower),
					Descriptor: model.Descriptor(line),
				}
				s.Tools = append(s.Tools, tool)
				tools = append(tools, model.ToolRef(tool.ID))
			}
			if strings.Contains(line, "AgentDefinition") {
				s.Agents = append(s.Agents, model.Agent{
					ID:        model.StableID("agent", "claude-sdk", path, line),
					Name:      pyutil.KeywordString(line, "name", "claude-agent"),
					Framework: model.FrameworkClaudeSDK,
					Location:  model.RelativeLocation(root, path, i+1),
					Tools:     append([]model.ToolRef{}, tools...),
					Metadata:  map[string]string{"detector": "claude_sdk"},
				})
			}
		}
		return nil
	})
}

func inferCapability(text string) model.Capability {
	cap := model.Capability{ReadsData: model.DataNone, ReadsUntrusted: model.UntrustedNone, Writes: model.WriteNone, Exfil: model.ExfilNone}
	if strings.Contains(text, "file") || strings.Contains(text, "secret") || strings.Contains(text, "github") {
		cap.ReadsData = model.DataSensitive
	}
	if strings.Contains(text, "web") || strings.Contains(text, "fetch") || strings.Contains(text, "mcp") {
		cap.ReadsUntrusted = model.UntrustedMCP
	}
	if strings.Contains(text, "slack") || strings.Contains(text, "chat") {
		cap.Writes = model.WriteExternal
		cap.Exfil = model.ExfilChat
	}
	return cap
}
