package generic

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/layergrid/layergrid-cli/internal/detectors/detectopts"
	"github.com/layergrid/layergrid-cli/internal/detectors/pyutil"
	"github.com/layergrid/layergrid-cli/internal/model"
)

type Detector struct{}

func New() Detector { return Detector{} }

func (Detector) Name() string               { return "generic" }
func (Detector) Framework() model.Framework { return model.FrameworkCustom }

func (Detector) Detect(root string, s *model.Stack, opts detectopts.Options) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", ".venv", "node_modules", "vendor", "dist", "build":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".py") {
			return nil
		}
		if !pyutil.Included(root, path, opts.Include, opts.Exclude) {
			return nil
		}
		return scanPython(root, path, s)
	})
}

func scanPython(root, path string, s *model.Stack) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	lineNo := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lineNo++
		line := scanner.Text()
		lower := strings.ToLower(line)
		if strings.Contains(lower, "openai.chat.completions.create") || strings.Contains(lower, "anthropic.messages.create") {
			s.Models = append(s.Models, model.Model{
				ID:       model.StableID("model", path, line),
				Provider: providerFor(lower),
				Name:     "unknown",
				Location: model.RelativeLocation(root, path, lineNo),
			})
		}
		if looksHardcodedKey(line) {
			s.Tools = append(s.Tools, evidenceTool(root, path, lineNo, "hardcoded-api-key", "hardcoded_key"))
		}
	}
	return scanner.Err()
}

func evidenceTool(root, path string, line int, name, marker string) model.Tool {
	return model.Tool{
		ID:         model.StableID("tool", "generic", path, name, marker),
		Name:       name,
		Kind:       model.ToolKindFunction,
		Source:     model.ToolSource{Kind: "python", Name: path},
		Location:   model.RelativeLocation(root, path, line),
		Capability: model.Capability{ReadsData: model.DataSensitive, ReadsUntrusted: model.UntrustedNone, Writes: model.WriteNone, Exfil: model.ExfilNone},
		Descriptor: model.Descriptor(name, marker),
		Metadata:   map[string]string{"evidence": marker},
	}
}

func providerFor(line string) string {
	if strings.Contains(line, "anthropic") {
		return "anthropic"
	}
	return "openai"
}

func looksHardcodedKey(line string) bool {
	lower := strings.ToLower(line)
	if strings.Contains(line, "sk-...") || strings.Contains(line, "sk-proj-1234567890") || strings.Contains(lower, "detector=") {
		return false
	}
	if strings.Contains(line, "ghp_XXXXXXXXXXXXXXXX") {
		return false
	}
	if !strings.Contains(lower, "api_key") && !strings.Contains(lower, "token") {
		return false
	}
	openAIKey := regexp.MustCompile(`['"]sk-[A-Za-z0-9_-]{20,}['"]`)
	githubKey := regexp.MustCompile(`['"]ghp_[A-Za-z0-9_]{20,}['"]`)
	return openAIKey.MatchString(line) || githubKey.MatchString(line)
}
