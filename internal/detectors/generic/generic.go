package generic

import (
	"bufio"
	"os"
	"path/filepath"
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
		if strings.Contains(lower, "os.environ") && (strings.Contains(lower, "token") || strings.Contains(lower, "key") || strings.Contains(lower, "secret")) {
			s.Tools = append(s.Tools, evidenceTool(root, path, lineNo, "env-credential-read", "env_credential_inline"))
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
	return (strings.Contains(lower, "api_key") || strings.Contains(lower, "token")) &&
		(strings.Contains(line, "\"sk-") || strings.Contains(line, "'sk-") || strings.Contains(line, "ghp_"))
}
