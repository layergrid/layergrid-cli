package openai_sdk

import (
	"strings"

	"github.com/layergrid/layergrid-cli/internal/detectors/detectopts"
	"github.com/layergrid/layergrid-cli/internal/detectors/pyutil"
	"github.com/layergrid/layergrid-cli/internal/model"
)

type Detector struct{}

func New() Detector { return Detector{} }

func (Detector) Name() string               { return "openai_sdk" }
func (Detector) Framework() model.Framework { return model.FrameworkOpenAISDK }

func (Detector) Detect(root string, s *model.Stack, opts detectopts.Options) error {
	return pyutil.Walk(root, opts.Include, opts.Exclude, func(path string, lines []string) error {
		if !pyutil.HasAny(lines, "from openai import OpenAI", "OpenAI(") {
			return nil
		}
		for i, line := range lines {
			lower := strings.ToLower(line)
			if strings.Contains(lower, "assistants.create") || strings.Contains(lower, "responses.create") {
				block := pyutil.Block(lines, i, 25)
				refs := emitTools(root, path, i+1, block, s)
				s.Agents = append(s.Agents, model.Agent{
					ID:        model.StableID("agent", "openai-sdk", path, line),
					Name:      pyutil.KeywordString(line, "name", "openai-agent"),
					Framework: model.FrameworkOpenAISDK,
					Location:  model.RelativeLocation(root, path, i+1),
					Tools:     append([]model.ToolRef{}, refs...),
					Metadata:  map[string]string{"detector": "openai_sdk"},
				})
			}
		}
		return nil
	})
}

func emitTools(root, path string, line int, block string, s *model.Stack) []model.ToolRef {
	lower := strings.ToLower(block)
	var refs []model.ToolRef
	for _, name := range []string{"code_interpreter", "file_search"} {
		if !strings.Contains(lower, name) {
			continue
		}
		kind := model.ToolKindFunction
		cap := model.Capability{ReadsData: model.DataInternal, ReadsUntrusted: model.UntrustedNone, Writes: model.WriteNone, Exfil: model.ExfilNone}
		if name == "code_interpreter" {
			kind = model.ToolKindCode
			cap = model.Capability{ReadsData: model.DataInternal, ReadsUntrusted: model.UntrustedNone, Writes: model.WriteLocal, Exfil: model.ExfilShell}
		}
		tool := model.Tool{
			ID:         model.StableID("tool", "openai-sdk", path, name, block),
			Name:       name,
			Kind:       kind,
			Source:     model.ToolSource{Kind: "python", Name: path},
			Location:   model.RelativeLocation(root, path, line),
			Capability: cap,
			Descriptor: model.Descriptor(block),
		}
		s.Tools = append(s.Tools, tool)
		refs = append(refs, model.ToolRef(tool.ID))
	}
	return refs
}
