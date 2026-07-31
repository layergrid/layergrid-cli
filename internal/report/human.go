package report

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/layergrid/layergrid-cli/internal/model"
	"github.com/layergrid/layergrid-cli/internal/scan"
	"github.com/layergrid/layergrid-cli/internal/trifecta"
)

type Human struct {
	NoColor bool
}

func (h Human) Format(r scan.Result) ([]byte, error) {
	color.NoColor = h.NoColor
	var b bytes.Buffer
	fmt.Fprintf(&b, "Scanning agent stack in %s...\n\n", r.Stack.Root)
	fmt.Fprintf(&b, "Discovered\n")
	fmt.Fprintf(&b, "  Agents      %-5d %s\n", len(r.Stack.Agents), frameworkSummary(r.Stack.Agents))
	fmt.Fprintf(&b, "  Tools       %-5d %s\n", len(r.Stack.Tools), toolSummary(r.Stack.Tools))
	fmt.Fprintf(&b, "  MCP Servers %-5d %s\n", len(r.Stack.MCPServers), mcpSummary(r.Stack.MCPServers))
	fmt.Fprintf(&b, "  Datasources %d\n\n", len(r.Stack.Datasources))

	criticalPaths := countSeverity(r.Findings, trifecta.SeverityCritical) + countSeverity(r.Findings, trifecta.SeverityHigh)
	if criticalPaths > 0 {
		fmt.Fprintf(&b, "Lethal Trifecta Signals  ·  %d findings\n\n", criticalPaths)
	} else {
		fmt.Fprintf(&b, "No lethal trifecta path detected\n\n")
	}
	for i, f := range firstFindings(r.Findings, 3) {
		fmt.Fprintf(&b, "  path #%d  %-8s score %+d    [%s]\n", i+1, strings.ToUpper(string(f.Severity)), f.ScoreImpact, f.RuleID)
		fmt.Fprintf(&b, "    %s\n", pathLine(f))
		fmt.Fprintf(&b, "    Fix: %s", strings.TrimSpace(f.Fix))
		if f.Location.Path != "" {
			fmt.Fprintf(&b, "  (see %s:%d)", f.Location.Path, f.Location.Line)
		}
		fmt.Fprintf(&b, "\n\n")
	}
	if len(r.Findings) > 3 {
		fmt.Fprintf(&b, "Other findings                    (%d)\n", len(r.Findings)-3)
		for _, f := range r.Findings[3:] {
			fmt.Fprintf(&b, "  %-7s %-28s %s\n", strings.ToUpper(string(f.Severity)), f.RuleID, f.Rationale)
		}
		fmt.Fprintf(&b, "\n")
	}
	fmt.Fprintf(&b, "----------------------------------------\n")
	fmt.Fprintf(&b, "  Trifecta Score      %d / 100\n", r.Score.Value)
	fmt.Fprintf(&b, "  Grade                %s\n", r.Score.Grade)
	fmt.Fprintf(&b, "  Findings             %d  (%d critical, %d high, %d medium, %d low)\n",
		len(r.Findings),
		r.Score.Counts[trifecta.SeverityCritical],
		r.Score.Counts[trifecta.SeverityHigh],
		r.Score.Counts[trifecta.SeverityMedium],
		r.Score.Counts[trifecta.SeverityLow],
	)
	fmt.Fprintf(&b, "  Scan time           %.1fs\n", r.Duration.Seconds())
	fmt.Fprintf(&b, "----------------------------------------\n\n")
	fmt.Fprintf(&b, "Run  layergrid explain LG-LETHAL-TRIFECTA-01  for details.\n")
	return b.Bytes(), nil
}

func firstFindings(findings []trifecta.Finding, n int) []trifecta.Finding {
	if len(findings) < n {
		return findings
	}
	return findings[:n]
}

func pathLine(f trifecta.Finding) string {
	if len(f.Path) == 0 {
		return fmt.Sprintf("%s  %s", f.Subject.Name, f.Rationale)
	}
	parts := make([]string, 0, len(f.Path))
	for _, node := range f.Path {
		if node.Name != "" {
			parts = append(parts, node.Name)
		} else {
			parts = append(parts, node.ID)
		}
	}
	return strings.Join(parts, " -> ")
}

func countSeverity(findings []trifecta.Finding, sev trifecta.Severity) int {
	count := 0
	for _, f := range findings {
		if f.Severity == sev {
			count++
		}
	}
	return count
}

func frameworkSummary(agents []model.Agent) string {
	if len(agents) == 0 {
		return ""
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
		parts = append(parts, fmt.Sprintf("%s x %d", key, counts[model.Framework(key)]))
	}
	return "(" + strings.Join(parts, ", ") + ")"
}

func toolSummary(tools []model.Tool) string {
	if len(tools) == 0 {
		return ""
	}
	counts := map[model.ToolKind]int{}
	for _, tool := range tools {
		counts[tool.Kind]++
	}
	return fmt.Sprintf("(%d function, %d MCP, %d shell, %d code)", counts[model.ToolKindFunction], counts[model.ToolKindMCP], counts[model.ToolKindShell], counts[model.ToolKindCode])
}

func mcpSummary(servers []model.MCPServer) string {
	if len(servers) == 0 {
		return ""
	}
	remote := 0
	unverified := 0
	for _, server := range servers {
		if server.IsExternal {
			remote++
		}
		if server.Publisher == "unknown" || server.Publisher == "" {
			unverified++
		}
	}
	return fmt.Sprintf("(%d remote, %d unverified publishers)", remote, unverified)
}
