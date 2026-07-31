package langchain

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/layergrid/layergrid-cli/internal/model"
)

type Detector struct{}

func New() Detector { return Detector{} }

func (Detector) Name() string               { return "langchain" }
func (Detector) Framework() model.Framework { return model.FrameworkLangChain }

func (d Detector) Detect(root string, s *model.Stack) error {
	return walkPython(root, func(path string) error {
		return d.detectFile(root, path, s)
	})
}

func (Detector) detectFile(root, path string, s *model.Stack) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	var lines []string
	hasLangChain := false
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		if strings.Contains(line, "langchain") || strings.Contains(line, "langgraph") {
			hasLangChain = true
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if !hasLangChain {
		return nil
	}

	toolVarByName := map[string]string{}
	decoratorTool := false
	nameRe := regexp.MustCompile(`name\s*=\s*["']([^"']+)["']`)
	assignRe := regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_]*)\s*=`)
	for i, line := range lines {
		lineNo := i + 1
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "@tool") {
			decoratorTool = true
			continue
		}
		if decoratorTool && strings.HasPrefix(trimmed, "def ") {
			name := strings.TrimSuffix(strings.TrimPrefix(strings.Fields(trimmed)[1], ""), "(")
			if idx := strings.Index(name, "("); idx >= 0 {
				name = name[:idx]
			}
			tool := makeTool(root, path, lineNo, name, localFunctionBlock(lines, i))
			s.Tools = append(s.Tools, tool)
			toolVarByName[name] = tool.ID
			decoratorTool = false
		}
		if strings.Contains(line, "Tool(") || strings.Contains(line, "StructuredTool.from_function") || strings.Contains(line, "FunctionTool.from_defaults") {
			name := "langchain-tool"
			if m := nameRe.FindStringSubmatch(line); len(m) == 2 {
				name = m[1]
			} else if m := assignRe.FindStringSubmatch(line); len(m) == 2 {
				name = m[1]
			}
			tool := makeTool(root, path, lineNo, name, localFunctionBlock(lines, i))
			s.Tools = append(s.Tools, tool)
			toolVarByName[name] = tool.ID
		}
		if isAgentConstruction(trimmed) {
			agent := model.Agent{
				ID:        model.StableID("agent", path, trimmed),
				Name:      agentName(line),
				Framework: model.FrameworkLangChain,
				Location:  model.RelativeLocation(root, path, lineNo),
				Tools:     refs(toolVarByName),
				Memory:    inferMemory(lines),
				Metadata:  map[string]string{"detector": "langchain"},
			}
			s.Agents = append(s.Agents, agent)
		}
	}
	return nil
}

func makeTool(root, path string, line int, name string, block string) model.Tool {
	cap := model.Capability{ReadsData: model.DataNone, ReadsUntrusted: model.UntrustedNone, Writes: model.WriteNone, Exfil: model.ExfilNone}
	lower := strings.ToLower(name + "\n" + block)
	switch {
	case strings.Contains(lower, "secret") || strings.Contains(lower, ".env") || strings.Contains(lower, "credential") || strings.Contains(lower, "api_key") || strings.Contains(lower, "token"):
		cap.ReadsData = model.DataSensitive
	case strings.Contains(lower, "pii") || strings.Contains(lower, "customer") || strings.Contains(lower, "database") || strings.Contains(lower, "db"):
		cap.ReadsData = model.DataSensitive
	}
	if strings.Contains(lower, "requests.") || strings.Contains(lower, "httpx.") || strings.Contains(lower, "urllib.") || strings.Contains(lower, "web") || strings.Contains(lower, "url") {
		cap.ReadsUntrusted = model.UntrustedWeb
		cap.NetworkEgress = []string{"*"}
	}
	if strings.Contains(lower, "rag") || strings.Contains(lower, "retriev") || strings.Contains(lower, "vector") || strings.Contains(lower, "upload") {
		cap.ReadsUntrusted = model.UntrustedRAG
	}
	switch {
	case strings.Contains(lower, "email") || strings.Contains(lower, "gmail"):
		cap.Writes = model.WriteExternal
		cap.Exfil = model.ExfilEmail
	case strings.Contains(lower, "slack") || strings.Contains(lower, "discord") || strings.Contains(lower, "teams") || strings.Contains(lower, "chat"):
		cap.Writes = model.WriteExternal
		cap.Exfil = model.ExfilChat
	case strings.Contains(lower, "post") || strings.Contains(lower, "send") || strings.Contains(lower, "webhook"):
		cap.Writes = model.WriteExternal
		cap.Exfil = model.ExfilHTTP
	case strings.Contains(lower, "subprocess") || strings.Contains(lower, "os.system"):
		cap.Writes = model.WriteExternal
		cap.Exfil = model.ExfilShell
	}
	kind := model.ToolKindFunction
	if cap.Exfil == model.ExfilShell {
		kind = model.ToolKindShell
	}
	return model.Tool{
		ID:         model.StableID("tool", path, name, block),
		Name:       name,
		Kind:       kind,
		Source:     model.ToolSource{Kind: "python", Name: path},
		Location:   model.RelativeLocation(root, path, line),
		Capability: cap,
		Descriptor: model.Descriptor(name, block),
		Metadata:   evidence(lower),
	}
}

func evidence(lower string) map[string]string {
	m := map[string]string{}
	if strings.Contains(lower, "os.environ") && (strings.Contains(lower, "token") || strings.Contains(lower, "key") || strings.Contains(lower, "secret")) {
		m["evidence"] = "env_credential_inline"
	}
	if strings.Contains(lower, "\"sk-") || strings.Contains(lower, "'sk-") || strings.Contains(lower, "ghp_") {
		m["evidence"] = "hardcoded_key"
	}
	if len(m) == 0 {
		return nil
	}
	return m
}

func refs(m map[string]string) []model.ToolRef {
	ids := make([]string, 0, len(m))
	for _, id := range m {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := make([]model.ToolRef, 0, len(ids))
	for _, id := range ids {
		out = append(out, model.ToolRef(id))
	}
	return out
}

func isAgentConstruction(line string) bool {
	if strings.HasPrefix(line, "from ") || strings.HasPrefix(line, "import ") {
		return false
	}
	return strings.Contains(line, "create_react_agent(") ||
		strings.Contains(line, "initialize_agent(") ||
		strings.Contains(line, "AgentExecutor(") ||
		strings.Contains(line, "create_openai_functions_agent(")
}

func agentName(line string) string {
	re := regexp.MustCompile(`name\s*=\s*["']([^"']+)["']`)
	if m := re.FindStringSubmatch(line); len(m) == 2 {
		return m[1]
	}
	return "langchain-agent"
}

func inferMemory(lines []string) model.MemoryConfig {
	joined := strings.ToLower(strings.Join(lines, "\n"))
	return model.MemoryConfig{Persistent: strings.Contains(joined, "memory"), Backend: ""}
}

func localFunctionBlock(lines []string, start int) string {
	end := start + 12
	if end > len(lines) {
		end = len(lines)
	}
	return strings.Join(lines[start:end], "\n")
}

func walkPython(root string, visit func(path string) error) error {
	return filepathWalk(root, ".py", visit)
}

func filepathWalk(root, suffix string, visit func(path string) error) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", ".venv", "node_modules", "vendor":
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, suffix) {
			return visit(path)
		}
		return nil
	})
}
