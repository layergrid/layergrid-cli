package llamaindex

import (
	"strings"

	"github.com/layergrid/layergrid-cli/internal/detectors/pyutil"
	"github.com/layergrid/layergrid-cli/internal/model"
)

type Detector struct{}

func New() Detector { return Detector{} }

func (Detector) Name() string               { return "llamaindex" }
func (Detector) Framework() model.Framework { return model.FrameworkLlamaIndex }

func (Detector) Detect(root string, s *model.Stack) error {
	return pyutil.Walk(root, func(path string, lines []string) error {
		if !pyutil.HasAny(lines, "llama_index", "llamaindex") {
			return nil
		}
		var tools []model.ToolRef
		for i, line := range lines {
			if strings.Contains(line, "VectorStoreIndex") || strings.Contains(line, "QueryEngine") {
				s.Datasources = append(s.Datasources, model.Datasource{
					ID:          model.StableID("datasource", "llamaindex", path, line),
					Kind:        "vector_store",
					Sensitivity: model.DataInternal,
					Trust:       trust(line),
					Location:    model.RelativeLocation(root, path, i+1),
					Description: "LlamaIndex retrieval source",
				})
				tool := model.Tool{
					ID:         model.StableID("tool", "llamaindex-rag", path, line),
					Name:       pyutil.AssignmentName(line, "llamaindex-rag"),
					Kind:       model.ToolKindFunction,
					Source:     model.ToolSource{Kind: "python", Name: path},
					Location:   model.RelativeLocation(root, path, i+1),
					Capability: model.Capability{ReadsData: model.DataInternal, ReadsUntrusted: model.UntrustedRAG, Writes: model.WriteNone, Exfil: model.ExfilNone},
					Descriptor: model.Descriptor(line),
				}
				s.Tools = append(s.Tools, tool)
				tools = append(tools, model.ToolRef(tool.ID))
			}
			if strings.Contains(line, "FunctionTool.from_defaults") {
				tool := model.Tool{
					ID:         model.StableID("tool", "llamaindex", path, line),
					Name:       pyutil.AssignmentName(line, "llamaindex-tool"),
					Kind:       model.ToolKindFunction,
					Source:     model.ToolSource{Kind: "python", Name: path},
					Location:   model.RelativeLocation(root, path, i+1),
					Capability: inferCapability(line + "\n" + pyutil.Block(lines, i, 8)),
					Descriptor: model.Descriptor(line),
				}
				s.Tools = append(s.Tools, tool)
				tools = append(tools, model.ToolRef(tool.ID))
			}
			if strings.Contains(line, "ReActAgent.from_tools") || strings.Contains(line, "AgentRunner") {
				s.Agents = append(s.Agents, model.Agent{
					ID:        model.StableID("agent", "llamaindex", path, line),
					Name:      "llamaindex-agent",
					Framework: model.FrameworkLlamaIndex,
					Location:  model.RelativeLocation(root, path, i+1),
					Tools:     append([]model.ToolRef{}, tools...),
					Metadata:  map[string]string{"detector": "llamaindex"},
				})
			}
		}
		return nil
	})
}

func trust(line string) string {
	if strings.Contains(strings.ToLower(line), "upload") || strings.Contains(strings.ToLower(line), "user") {
		return "user_supplied"
	}
	return "external"
}

func inferCapability(text string) model.Capability {
	lower := strings.ToLower(text)
	cap := model.Capability{ReadsData: model.DataNone, ReadsUntrusted: model.UntrustedNone, Writes: model.WriteNone, Exfil: model.ExfilNone}
	if strings.Contains(lower, "secret") || strings.Contains(lower, "pii") || strings.Contains(lower, "customer") {
		cap.ReadsData = model.DataSensitive
	}
	if strings.Contains(lower, "email") || strings.Contains(lower, "gmail") {
		cap.Writes = model.WriteExternal
		cap.Exfil = model.ExfilEmail
	}
	return cap
}
