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
			if strings.Contains(lower, "code_interpreter") || strings.Contains(lower, "file_search") {
				kind := model.ToolKindFunction
				cap := model.Capability{ReadsData: model.DataInternal, ReadsUntrusted: model.UntrustedNone, Writes: model.WriteNone, Exfil: model.ExfilNone}
				name := "file_search"
				if strings.Contains(lower, "code_interpreter") {
					kind = model.ToolKindCode
					cap = model.Capability{ReadsData: model.DataInternal, ReadsUntrusted: model.UntrustedNone, Writes: model.WriteLocal, Exfil: model.ExfilShell}
					name = "code_interpreter"
				}
				tool := model.Tool{
					ID:         model.StableID("tool", "openai-sdk", path, name, line),
					Name:       name,
					Kind:       kind,
					Source:     model.ToolSource{Kind: "python", Name: path},
					Location:   model.RelativeLocation(root, path, i+1),
					Capability: cap,
					Descriptor: model.Descriptor(line),
				}
				s.Tools = append(s.Tools, tool)
			}
		}
		var refs []model.ToolRef
		for _, tool := range s.Tools {
			if tool.Source.Name == path {
				refs = append(refs, model.ToolRef(tool.ID))
			}
		}
		for i, line := range lines {
			lower := strings.ToLower(line)
			if strings.Contains(lower, "assistants.create") || strings.Contains(lower, "responses.create") {
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
