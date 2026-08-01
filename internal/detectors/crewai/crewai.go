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
		if pyutil.FrameworkSource(root, path, "crewai") {
			return nil
		}
		if !pyutil.HasAny(lines, "from crewai import", "crewai") {
			return nil
		}
		toolIDsByName := map[string]model.ToolRef{}
		inDocstring := false
		for i, line := range lines {
			if togglesDocstring(line) {
				inDocstring = !inDocstring
				continue
			}
			if inDocstring {
				continue
			}
			lower := strings.ToLower(line)
			if strings.Contains(line, "Tool(") || strings.Contains(lower, "tools.") {
				name := pyutil.AssignmentName(line, "crewai-tool")
				tool := inferTool(root, path, i+1, name, pyutil.Block(lines, i, 10))
				s.Tools = append(s.Tools, tool)
				toolIDsByName[name] = model.ToolRef(tool.ID)
			}
			if strings.Contains(line, "Agent(") {
				block := pyutil.Block(lines, i, 25)
				agent := model.Agent{
					ID:        model.StableID("agent", "crewai", path, line),
					Name:      pyutil.KeywordString(line, "role", "crewai-agent"),
					Framework: model.FrameworkCrewAI,
					Location:  model.RelativeLocation(root, path, i+1),
					Tools:     refsInBlock(toolIDsByName, block),
					Memory:    memory(block),
					Metadata:  map[string]string{"detector": "crewai"},
				}
				s.Agents = append(s.Agents, agent)
			}
		}
		return nil
	})
}

func refsInBlock(toolIDsByName map[string]model.ToolRef, block string) []model.ToolRef {
	if len(toolIDsByName) == 0 || !strings.Contains(strings.ToLower(block), "tools") {
		return nil
	}
	var refs []model.ToolRef
	for name, ref := range toolIDsByName {
		if strings.Contains(block, name) {
			refs = append(refs, ref)
		}
	}
	return refs
}

func togglesDocstring(line string) bool {
	return strings.Count(line, `"""`)%2 == 1 || strings.Count(line, `'''`)%2 == 1
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

func memory(line string) model.MemoryConfig {
	joined := strings.ToLower(line)
	return model.MemoryConfig{
		Persistent:        strings.Contains(joined, "memory=true") || strings.Contains(joined, "memory = true"),
		SharedAcrossUsers: strings.Contains(joined, "shared_memory") || strings.Contains(joined, "cross_user"),
	}
}
