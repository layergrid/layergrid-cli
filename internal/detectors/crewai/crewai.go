package crewai

import (
	"strings"

	"github.com/layergrid/layergrid-cli/internal/detectors/detectopts"
	"github.com/layergrid/layergrid-cli/internal/detectors/pyutil"
	"github.com/layergrid/layergrid-cli/internal/model"
)

type Detector struct{}

func New() Detector { return Detector{} }

func (Detector) Name() string               { return "crewai" }
func (Detector) Framework() model.Framework { return model.FrameworkCrewAI }

func (Detector) Detect(root string, s *model.Stack, opts detectopts.Options) error {
	return pyutil.Walk(root, opts.Include, opts.Exclude, func(path string, lines []string) error {
		if !pyutil.HasAny(lines, "from crewai import", "crewai") {
			return nil
		}
		var toolIDs []model.ToolRef
		for i, line := range lines {
			lower := strings.ToLower(line)
			if strings.Contains(line, "Tool(") || strings.Contains(lower, "tools.") {
				tool := inferTool(root, path, i+1, pyutil.AssignmentName(line, "crewai-tool"), pyutil.Block(lines, i, 10))
				s.Tools = append(s.Tools, tool)
				toolIDs = append(toolIDs, model.ToolRef(tool.ID))
			}
			if strings.Contains(line, "Agent(") {
				agent := model.Agent{
					ID:        model.StableID("agent", "crewai", path, line),
					Name:      pyutil.KeywordString(line, "role", "crewai-agent"),
					Framework: model.FrameworkCrewAI,
					Location:  model.RelativeLocation(root, path, i+1),
					Tools:     append([]model.ToolRef{}, toolIDs...),
					Memory:    memory(lines),
					Metadata:  map[string]string{"detector": "crewai"},
				}
				s.Agents = append(s.Agents, agent)
			}
			if strings.Contains(line, "Crew(") && len(s.Agents) > 1 {
				parent := s.Agents[len(s.Agents)-1]
				for _, agent := range s.Agents[:len(s.Agents)-1] {
					if agent.Framework == model.FrameworkCrewAI {
						parent.SubAgents = append(parent.SubAgents, model.AgentRef(agent.ID))
					}
				}
				s.Agents[len(s.Agents)-1] = parent
			}
		}
		return nil
	})
}

func inferTool(root, path string, line int, name, block string) model.Tool {
	lower := strings.ToLower(name + "\n" + block)
	cap := model.Capability{ReadsData: model.DataNone, ReadsUntrusted: model.UntrustedNone, Writes: model.WriteNone, Exfil: model.ExfilNone}
	if strings.Contains(lower, "customer") || strings.Contains(lower, "pii") || strings.Contains(lower, "secret") || strings.Contains(lower, "credential") {
		cap.ReadsData = model.DataSensitive
	}
	if strings.Contains(lower, "rag") || strings.Contains(lower, "upload") || strings.Contains(lower, "ticket") || strings.Contains(lower, "inbox") {
		cap.ReadsUntrusted = model.UntrustedInbox
	}
	if strings.Contains(lower, "slack") || strings.Contains(lower, "chat") {
		cap.Writes = model.WriteExternal
		cap.Exfil = model.ExfilChat
	}
	if strings.Contains(lower, "email") || strings.Contains(lower, "gmail") {
		cap.Writes = model.WriteExternal
		cap.Exfil = model.ExfilEmail
	}
	return model.Tool{
		ID:         model.StableID("tool", "crewai", path, name, block),
		Name:       name,
		Kind:       model.ToolKindFunction,
		Source:     model.ToolSource{Kind: "python", Name: path},
		Location:   model.RelativeLocation(root, path, line),
		Capability: cap,
		Descriptor: model.Descriptor(name, block),
	}
}

func memory(lines []string) model.MemoryConfig {
	joined := strings.ToLower(strings.Join(lines, "\n"))
	return model.MemoryConfig{
		Persistent:        strings.Contains(joined, "memory=true") || strings.Contains(joined, "memory = true"),
		SharedAcrossUsers: strings.Contains(joined, "shared_memory") || strings.Contains(joined, "cross_user"),
	}
}
