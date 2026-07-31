package report

import (
	"bytes"
	"fmt"
	"html"
	"strings"

	"github.com/layergrid/layergrid-cli/internal/scan"
	"github.com/layergrid/layergrid-cli/internal/trifecta"
)

type HTML struct{}

func (HTML) Format(r scan.Result) ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "<!doctype html><html><head><meta charset=\"utf-8\"><title>LayerGrid Scan</title>")
	fmt.Fprintf(&b, "<style>body{font-family:ui-sans-serif,system-ui;margin:0;color:#172026;background:#f7f8fa}main{max-width:1120px;margin:auto;padding:32px}header{display:flex;justify-content:space-between;align-items:end;border-bottom:1px solid #d9dee5;padding-bottom:20px}.score{font-size:48px;font-weight:750}.grade{font-size:28px}.band{background:#fff;border:1px solid #d9dee5;border-radius:8px;padding:20px;margin:20px 0}table{border-collapse:collapse;width:100%%;background:#fff}td,th{border-bottom:1px solid #e4e8ee;padding:10px;text-align:left;vertical-align:top}code{background:#eef1f5;padding:2px 5px;border-radius:4px}.pill{font-size:12px;text-transform:uppercase;font-weight:700}.critical,.high{color:#b42318}.medium{color:#b54708}.low{color:#475467}.graph{display:flex;gap:8px;flex-wrap:wrap}.node{border:1px solid #ccd3dd;border-radius:8px;padding:8px 10px;background:#fff}.arrow{align-self:center;color:#697586}</style></head><body><main>")
	fmt.Fprintf(&b, "<header><div><h1>LayerGrid Scan</h1><p>%s</p></div><div><div class=\"score\">%d / 100</div><div class=\"grade\">Grade %s</div></div></header>", html.EscapeString(r.Stack.Root), r.Score.Value, html.EscapeString(r.Score.Grade))
	fmt.Fprintf(&b, "<section class=\"band\"><h2>Discovered</h2><p>%d agents · %d tools · %d MCP servers · %d datasources</p></section>", len(r.Stack.Agents), len(r.Stack.Tools), len(r.Stack.MCPServers), len(r.Stack.Datasources))
	if len(r.Findings) > 0 {
		fmt.Fprintf(&b, "<section class=\"band\"><h2>Primary Path</h2><div class=\"graph\">%s</div></section>", graphHTML(r.Findings[0]))
	}
	fmt.Fprintf(&b, "<section class=\"band\"><h2>Findings</h2><table><tr><th>Severity</th><th>Rule</th><th>Location</th><th>Rationale</th><th>Fix</th></tr>")
	for _, f := range r.Findings {
		fmt.Fprintf(&b, "<tr><td class=\"pill %s\">%s</td><td><code>%s</code><br>%s</td><td><code>%s:%d</code></td><td>%s</td><td>%s</td></tr>", html.EscapeString(string(f.Severity)), html.EscapeString(string(f.Severity)), html.EscapeString(f.RuleID), html.EscapeString(f.RuleName), html.EscapeString(f.Location.Path), f.Location.Line, html.EscapeString(f.Rationale), html.EscapeString(f.Fix))
	}
	fmt.Fprintf(&b, "</table></section><section class=\"band\"><h2>Share</h2><p>Coming in v0.2</p></section></main></body></html>")
	return b.Bytes(), nil
}

func graphHTML(f trifecta.Finding) string {
	if len(f.Path) == 0 {
		return fmt.Sprintf("<div class=\"node\">%s</div>", html.EscapeString(f.Subject.Name))
	}
	var parts []string
	for _, node := range f.Path {
		label := node.Name
		if label == "" {
			label = node.ID
		}
		parts = append(parts, fmt.Sprintf("<div class=\"node\"><strong>%s</strong><br><small>%s</small></div>", html.EscapeString(label), html.EscapeString(node.Kind)))
	}
	return strings.Join(parts, "<div class=\"arrow\">-&gt;</div>")
}
