package report

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/layergrid/layergrid-cli/internal/model"
	"github.com/layergrid/layergrid-cli/internal/scan"
	"github.com/layergrid/layergrid-cli/internal/trifecta"
	"github.com/layergrid/layergrid-cli/internal/version"
)

const (
	brandOrange = "#FF4D00"
	designWidth = 100
)

type Human struct {
	NoColor bool
	Width   int
}

func (h Human) Format(r scan.Result) ([]byte, error) {
	renderer := humanRenderer{noColor: noColorEnabled(h.NoColor), width: resolveWidth(h.Width)}
	return []byte(renderer.scan(r)), nil
}

type humanRenderer struct {
	noColor bool
	width   int
}

func (r humanRenderer) scan(result scan.Result) string {
	var b bytes.Buffer
	b.WriteString(r.scanHeader(result))
	b.WriteString("\n\n\n")
	b.WriteString(r.sectionTitle("Discovered", ""))
	b.WriteString("\n")
	b.WriteString(indent(r.discoveredTable(result.Stack), 2))
	b.WriteString("\n\n\n")
	if len(result.Findings) == 0 {
		b.WriteString(r.sectionTitle("Findings", "0 total"))
		b.WriteString("\n")
		b.WriteString(indent(r.emptyState(result.Score), 2))
		b.WriteString("\n\n\n")
	} else {
		b.WriteString(r.sectionTitle("Findings", fmt.Sprintf("%d total", len(result.Findings))))
		b.WriteString("\n")
		b.WriteString(indent(r.findingsTable(result.Findings), 2))
		b.WriteString("\n\n\n")
		attack := highImpactFindings(result.Findings)
		if len(attack) > 0 {
			b.WriteString(r.sectionTitle("Attack Paths", "showing CRITICAL and HIGH"))
			b.WriteString("\n")
			b.WriteString(indent(r.attackTable(attack), 2))
			b.WriteString("\n\n\n")
		}
	}
	b.WriteString(r.scoreBlock(result))
	b.WriteString("\n\n\n")
	b.WriteString("  Next  layergrid explain LG-LETHAL-TRIFECTA-01     (deep-dive)\n")
	b.WriteString("        layergrid scan --format html -o report.html (shareable)\n")
	b.WriteString("        layergrid scan --fail-on high                (CI mode)")
	return strings.TrimRight(b.String(), "\n") + "\n"
}

func (r humanRenderer) scanHeader(result scan.Result) string {
	files := countFiles(result.Stack.Root)
	line := fmt.Sprintf("Scanning %s\nReading %d files  %s  100%%  %.1fs", result.Stack.Root, files, progressBar(20, 1), result.Duration.Seconds())
	return lipgloss.NewStyle().
		Width(minInt(designWidth-4, maxInt(56, r.width-4))).
		Padding(1, 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(r.color(brandOrange)).
		Render(line)
}

func (r humanRenderer) sectionTitle(left, right string) string {
	if right == "" {
		return "  " + r.bold(left)
	}
	w := minInt(designWidth-4, maxInt(56, r.width-4))
	gap := maxInt(1, w-lipgloss.Width(left)-lipgloss.Width(right))
	return "  " + r.bold(left) + strings.Repeat(" ", gap) + right
}

func (r humanRenderer) discoveredTable(stack model.Stack) string {
	rows := [][]string{
		{"Agents", fmt.Sprintf("%d", len(stack.Agents)), frameworkSummary(stack.Agents)},
		{"Tools", fmt.Sprintf("%d", len(stack.Tools)), toolSummary(stack.Tools)},
		{"MCP Servers", fmt.Sprintf("%d", len(stack.MCPServers)), mcpSummary(stack.MCPServers)},
		{"Datasources", fmt.Sprintf("%d", len(stack.Datasources)), datasourceSummary(stack.Datasources)},
	}
	return r.table([]string{"Type", "Count", "Detail"}, rows, 78, nil)
}

func (r humanRenderer) findingsTable(findings []trifecta.Finding) string {
	rows := make([][]string, 0, len(findings))
	for i, f := range findings {
		rows = append(rows, []string{
			fmt.Sprintf("%d", i+1),
			strings.ToUpper(string(f.Severity)),
			truncate(f.RuleID, 30),
			scoreText(f.ScoreImpact),
			truncate(locationText(f), 24),
		})
	}
	return r.table([]string{"#", "Severity", "Rule", "Score", "Location"}, rows, 92, func(row, col int) lipgloss.Style {
		s := baseCellStyle(row, col)
		if row == table.HeaderRow {
			return s.Bold(true).Align(lipgloss.Center)
		}
		if col == 0 {
			s = s.Align(lipgloss.Right)
		}
		if col == 1 || col == 3 {
			s = s.Align(lipgloss.Center)
			if row > 0 && row-1 < len(findings) {
				s = s.Foreground(r.severityColor(findings[row-1].Severity))
				if findings[row-1].Severity == trifecta.SeverityCritical || findings[row-1].Severity == trifecta.SeverityHigh {
					s = s.Bold(true)
				}
			}
		}
		return s
	})
}

func (r humanRenderer) attackTable(findings []trifecta.Finding) string {
	rows := make([][]string, 0, len(findings))
	for i, f := range findings {
		rows = append(rows, []string{
			fmt.Sprintf("%d", i+1),
			strings.ToUpper(string(f.Severity)),
			truncateMultiline(pathText(f), 38),
			truncateMultiline(fixText(f), 32),
		})
	}
	return table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(r.subtleBorder()).
		BorderRow(true).
		Width(92).
		StyleFunc(func(row, col int) lipgloss.Style {
			s := baseCellStyle(row, col)
			if row == table.HeaderRow {
				return s.Bold(true).Align(lipgloss.Center)
			}
			if col == 0 || col == 1 {
				s = s.Align(lipgloss.Center)
			}
			if col == 1 && row > 0 && row-1 < len(findings) {
				s = s.Foreground(r.severityColor(findings[row-1].Severity)).Bold(true)
			}
			return s
		}).
		Headers("#", "Severity", "Path", "Fix").
		Rows(rows...).
		String()
}

func (r humanRenderer) scoreBlock(result scan.Result) string {
	score := r.plainTable([][]string{
		{"Trifecta Score", fmt.Sprintf("%d / 100", result.Score.Value)},
		{"Grade", result.Score.Grade},
		{"Rubric", "v" + version.RubricVersion},
		{"Scan time", fmt.Sprintf("%.1fs", result.Duration.Seconds())},
	}, 34, func(row, col int) lipgloss.Style {
		s := lipgloss.NewStyle().Padding(0, 1)
		if col == 1 && row == 1 {
			s = s.Foreground(r.gradeColor(result.Score.Grade)).Bold(true).Align(lipgloss.Center)
		}
		return s
	})
	breakdown := r.plainTable([][]string{
		{"Critical", fmt.Sprintf("%d", result.Score.Counts[trifecta.SeverityCritical])},
		{"High", fmt.Sprintf("%d", result.Score.Counts[trifecta.SeverityHigh])},
		{"Medium", fmt.Sprintf("%d", result.Score.Counts[trifecta.SeverityMedium])},
		{"Low", fmt.Sprintf("%d", result.Score.Counts[trifecta.SeverityLow])},
	}, 22, nil)
	if r.width >= designWidth {
		return "  " + r.bold("Score") + strings.Repeat(" ", 46) + r.bold("Breakdown") + "\n" + lipgloss.JoinHorizontal(lipgloss.Top, indent(score, 2), strings.Repeat(" ", 18), breakdown)
	}
	return "  " + r.bold("Score") + "\n" + indent(score, 2) + "\n\n  " + r.bold("Breakdown") + "\n" + indent(breakdown, 2)
}

func (r humanRenderer) emptyState(score trifecta.Score) string {
	text := fmt.Sprintf("✓  No findings — Grade %s (%d/100).", score.Grade, score.Value)
	if r.noColor {
		text = fmt.Sprintf("✓  No findings — Grade %s (%d/100).", score.Grade, score.Value)
	}
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(r.subtleColor()).
		Padding(0, 2).
		Render(text)
}

func (r humanRenderer) table(headers []string, rows [][]string, width int, style table.StyleFunc) string {
	if style == nil {
		style = func(row, col int) lipgloss.Style {
			s := baseCellStyle(row, col)
			if row == table.HeaderRow {
				return s.Bold(true).Align(lipgloss.Center)
			}
			if col == 1 {
				return s.Align(lipgloss.Right)
			}
			return s
		}
	}
	return table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(r.subtleBorder()).
		Width(minInt(width, maxInt(40, r.width-6))).
		StyleFunc(style).
		Headers(headers...).
		Rows(rows...).
		String()
}

func (r humanRenderer) plainTable(rows [][]string, width int, style table.StyleFunc) string {
	if style == nil {
		style = func(row, col int) lipgloss.Style {
			s := lipgloss.NewStyle().Padding(0, 1)
			if col == 1 {
				return s.Align(lipgloss.Right)
			}
			return s
		}
	}
	return table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(r.subtleBorder()).
		Width(minInt(width, maxInt(20, r.width-6))).
		StyleFunc(style).
		Rows(rows...).
		String()
}

func baseCellStyle(row, col int) lipgloss.Style {
	s := lipgloss.NewStyle().Padding(0, 1)
	if row == table.HeaderRow {
		return s.Bold(true).Align(lipgloss.Center)
	}
	return s
}

func (r humanRenderer) color(hex string) lipgloss.TerminalColor {
	if r.noColor {
		return lipgloss.NoColor{}
	}
	return lipgloss.Color(hex)
}

func (r humanRenderer) subtleColor() lipgloss.TerminalColor {
	if r.noColor {
		return lipgloss.NoColor{}
	}
	return lipgloss.Color("240")
}

func (r humanRenderer) subtleBorder() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(r.subtleColor())
}

func (r humanRenderer) severityColor(sev trifecta.Severity) lipgloss.TerminalColor {
	if r.noColor {
		return lipgloss.NoColor{}
	}
	switch sev {
	case trifecta.SeverityCritical:
		return lipgloss.Color("196")
	case trifecta.SeverityHigh:
		return lipgloss.Color(brandOrange)
	case trifecta.SeverityMedium:
		return lipgloss.Color("220")
	default:
		return lipgloss.Color("245")
	}
}

func (r humanRenderer) gradeColor(grade string) lipgloss.TerminalColor {
	if r.noColor {
		return lipgloss.NoColor{}
	}
	switch grade {
	case "A", "B":
		return lipgloss.Color("42")
	case "C":
		return lipgloss.Color("220")
	case "D":
		return lipgloss.Color(brandOrange)
	default:
		return lipgloss.Color("196")
	}
}

func (r humanRenderer) bold(s string) string {
	if r.noColor {
		return s
	}
	return lipgloss.NewStyle().Bold(true).Render(s)
}

func highImpactFindings(findings []trifecta.Finding) []trifecta.Finding {
	out := make([]trifecta.Finding, 0, len(findings))
	for _, f := range findings {
		if f.Severity == trifecta.SeverityCritical || f.Severity == trifecta.SeverityHigh {
			out = append(out, f)
		}
	}
	return out
}

func frameworkSummary(agents []model.Agent) string {
	if len(agents) == 0 {
		return "—"
	}
	counts := map[model.Framework]int{}
	for _, agent := range agents {
		counts[agent.Framework]++
	}
	keys := make([]string, 0, len(counts))
	for fw := range counts {
		keys = append(keys, string(fw))
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s × %d", displayFramework(model.Framework(key)), counts[model.Framework(key)]))
	}
	return strings.Join(parts, ", ")
}

func toolSummary(tools []model.Tool) string {
	if len(tools) == 0 {
		return "—"
	}
	counts := map[model.ToolKind]int{}
	for _, tool := range tools {
		counts[tool.Kind]++
	}
	other := len(tools) - counts[model.ToolKindFunction] - counts[model.ToolKindMCP] - counts[model.ToolKindShell] - counts[model.ToolKindCode]
	if other < 0 {
		other = 0
	}
	return fmt.Sprintf("%d function, %d MCP, %d shell, %d code, %d other", counts[model.ToolKindFunction], counts[model.ToolKindMCP], counts[model.ToolKindShell], counts[model.ToolKindCode], other)
}

func mcpSummary(servers []model.MCPServer) string {
	if len(servers) == 0 {
		return "—"
	}
	remote, stdio, unverified := 0, 0, 0
	for _, server := range servers {
		if server.IsExternal {
			remote++
		}
		if server.Transport == "stdio" {
			stdio++
		}
		if server.Publisher == "unknown" || server.Publisher == "" {
			unverified++
		}
	}
	return fmt.Sprintf("%d stdio, %d remote · %d unverified publishers", stdio, remote, unverified)
}

func datasourceSummary(datasources []model.Datasource) string {
	if len(datasources) == 0 {
		return "—"
	}
	counts := map[model.DataAccess]int{}
	for _, ds := range datasources {
		counts[ds.Sensitivity]++
	}
	return fmt.Sprintf("%d sensitive, %d restricted, %d internal", counts[model.DataSensitive], counts[model.DataRestricted], counts[model.DataInternal])
}

func displayFramework(f model.Framework) string {
	switch f {
	case model.FrameworkLangChain:
		return "LangChain"
	case model.FrameworkLangGraph:
		return "LangGraph"
	case model.FrameworkLlamaIndex:
		return "LlamaIndex"
	case model.FrameworkCrewAI:
		return "CrewAI"
	case model.FrameworkAutoGen:
		return "AutoGen"
	case model.FrameworkClaudeSDK:
		return "Claude SDK"
	case model.FrameworkOpenAISDK:
		return "OpenAI SDK"
	default:
		return string(f)
	}
}

func pathText(f trifecta.Finding) string {
	nodes := f.Path
	if len(nodes) == 0 {
		name := f.Subject.Name
		if name == "" {
			name = f.Subject.ID
		}
		if name == "" {
			name = f.RuleID
		}
		return name
	}
	lines := make([]string, 0, len(nodes))
	for i, node := range nodes {
		name := node.Name
		if name == "" {
			name = node.ID
		}
		if i == 0 {
			lines = append(lines, name)
		} else {
			lines = append(lines, "→ "+name)
		}
	}
	return strings.Join(lines, "\n")
}

func fixText(f trifecta.Finding) string {
	fix := sentenceCase(strings.TrimSpace(f.Fix))
	if fix == "" {
		fix = "Review this capability composition."
	}
	loc := locationText(f)
	if loc == "—" {
		return fix
	}
	return fix + "\n\n" + loc
}

func locationText(f trifecta.Finding) string {
	if f.Location.Path == "" {
		if f.Subject.Name != "" {
			return f.Subject.Name
		}
		return "—"
	}
	if f.Location.Line > 0 {
		return fmt.Sprintf("%s:%d", f.Location.Path, f.Location.Line)
	}
	return f.Location.Path
}

func sentenceCase(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func scoreText(score int) string {
	if score < 0 {
		return fmt.Sprintf("−%d", -score)
	}
	return fmt.Sprintf("+%d", score)
}

func progressBar(width int, ratio float64) string {
	filled := int(float64(width) * ratio)
	if filled > width {
		filled = width
	}
	return strings.Repeat("▓", filled) + strings.Repeat("░", width-filled)
}

func countFiles(root string) int {
	if root == "" {
		return 0
	}
	count := 0
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", ".gocache", ".golangci-cache", "node_modules", ".venv":
				return filepath.SkipDir
			}
			return nil
		}
		count++
		return nil
	})
	return count
}

func noColorEnabled(explicit bool) bool {
	return explicit || os.Getenv("NO_COLOR") != ""
}

func resolveWidth(explicit int) int {
	if explicit > 0 {
		return explicit
	}
	if cols := os.Getenv("COLUMNS"); cols != "" {
		var n int
		if _, err := fmt.Sscanf(cols, "%d", &n); err == nil && n > 0 {
			return n
		}
	}
	return designWidth
}

func truncate(s string, max int) string {
	if lipgloss.Width(s) <= max {
		return s
	}
	var b strings.Builder
	for _, rr := range s {
		if lipgloss.Width(b.String()+string(rr)+"…") > max {
			break
		}
		b.WriteRune(rr)
	}
	return b.String() + "…"
}

func truncateMultiline(s string, max int) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = truncate(line, max)
	}
	return strings.Join(lines, "\n")
}

func indent(s string, spaces int) string {
	prefix := strings.Repeat(" ", spaces)
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var _ = time.Duration(0)
