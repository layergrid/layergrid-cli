package report

import (
	"bytes"
	"fmt"
	"html"

	"github.com/layergrid/layergrid/internal/scan"
)

type HTML struct{}

func (HTML) Format(r scan.Result) ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "<!doctype html><meta charset=\"utf-8\"><title>LayerGrid Scan</title>")
	fmt.Fprintf(&b, "<style>body{font-family:ui-sans-serif,system-ui;margin:32px;color:#172026}table{border-collapse:collapse;width:100%%}td,th{border-bottom:1px solid #ddd;padding:8px;text-align:left}.score{font-size:40px;font-weight:700}</style>")
	fmt.Fprintf(&b, "<h1>LayerGrid Scan</h1><div class=\"score\">%d / 100 · Grade %s</div>", r.Score.Value, html.EscapeString(r.Score.Grade))
	fmt.Fprintf(&b, "<h2>Findings</h2><table><tr><th>Severity</th><th>Rule</th><th>Location</th><th>Fix</th></tr>")
	for _, f := range r.Findings {
		fmt.Fprintf(&b, "<tr><td>%s</td><td>%s</td><td>%s:%d</td><td>%s</td></tr>", html.EscapeString(string(f.Severity)), html.EscapeString(f.RuleID), html.EscapeString(f.Location.Path), f.Location.Line, html.EscapeString(f.Fix))
	}
	fmt.Fprintf(&b, "</table>")
	return b.Bytes(), nil
}
