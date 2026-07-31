package generic

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/layergrid/layergrid/internal/model"
)

type Detector struct{}

func New() Detector { return Detector{} }

func (Detector) Name() string               { return "generic" }
func (Detector) Framework() model.Framework { return model.FrameworkCustom }

func (Detector) Detect(root string, s *model.Stack) error {
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
			s.Errors = append(s.Errors, model.ScanError{
				Detector: "generic",
				Message:  "possible hardcoded credential-like value; emitted as scan error pending precise rule support",
				Location: model.RelativeLocation(root, path, lineNo),
			})
		}
	}
	return scanner.Err()
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
