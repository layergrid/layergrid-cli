package autogen

import (
	"strings"

	"github.com/layergrid/layergrid-cli/internal/detectors/detectopts"
	"github.com/layergrid/layergrid-cli/internal/detectors/pyutil"
	"github.com/layergrid/layergrid-cli/internal/model"
)

type Detector struct{}

func New() Detector { return Detector{} }

func (Detector) Name() string               { return "autogen" }
func (Detector) Framework() model.Framework { return model.FrameworkAutoGen }

func (Detector) Detect(root string, s *model.Stack, opts detectopts.Options) error {
	return pyutil.Walk(root, opts.Include, opts.Exclude, func(path string, lines []string) error {
		if pyutil.FrameworkSource(root, path, "autogen") {
			return nil
		}
		if !pyutil.HasAny(lines, "autogen", "ConversableAgent", "GroupChat") {
			return nil
		}
		toolIDsByName := map[string]model.ToolRef{}
		for i, line := range lines {
			lower := strings.ToLower(line)
			if isExecutorConstruction(line) {
				name := pyutil.AssignmentName(line, "autogen-executor")
				kind := model.ToolKindCode
				rationale := "docker-code-exec"
				if strings.Contains(line, "LocalCommandLineCodeExecutor") {
					kind = model.ToolKindShell
					rationale = "local-shell-exec"
				}
				tool := model.Tool{
					ID:         model.StableID("tool", "autogen", path, line),
					Name:       name,
					Kind:       kind,
					Source:     model.ToolSource{Kind: "python", Name: path},
					Location:   model.RelativeLocation(root, path, i+1),
					Capability: model.Capability{ReadsData: model.DataNone, ReadsUntrusted: model.UntrustedNone, Writes: model.WriteExternal, Exfil: model.ExfilShell},
					Descriptor: model.Descriptor(line),
					Metadata:   map[string]string{"evidence": rationale},
				}
				s.Tools = append(s.Tools, tool)
				toolIDsByName[name] = model.ToolRef(tool.ID)
			}
			if strings.Contains(line, "ConversableAgent(") || strings.Contains(line, "AssistantAgent(") {
				block := pyutil.Block(lines, i, 25)
				s.Agents = append(s.Agents, model.Agent{
					ID:        model.StableID("agent", "autogen", path, line),
					Name:      pyutil.KeywordString(line, "name", "autogen-agent"),
					Framework: model.FrameworkAutoGen,
					Location:  model.RelativeLocation(root, path, i+1),
					Tools:     refsInBlock(toolIDsByName, block),
					Metadata:  map[string]string{"detector": "autogen", "source": lower},
				})
			}
		}
		return nil
	})
}

func refsInBlock(toolIDsByName map[string]model.ToolRef, block string) []model.ToolRef {
	if len(toolIDsByName) == 0 {
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

func isExecutorConstruction(line string) bool {
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "from ") || strings.HasPrefix(trimmed, "import ") {
		return false
	}
	return strings.Contains(line, "LocalCommandLineCodeExecutor(") || strings.Contains(line, "DockerCommandLineCodeExecutor(")
}
