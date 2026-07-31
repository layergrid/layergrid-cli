package report

import (
	"bytes"
	"fmt"

	"github.com/layergrid/layergrid-cli/internal/scan"
)

type SARIF struct{}

func (SARIF) Format(r scan.Result) ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprintf(&b, `{"version":"2.1.0","$schema":"https://json.schemastore.org/sarif-2.1.0.json","runs":[{"tool":{"driver":{"name":"LayerGrid","informationUri":"https://github.com/layergrid/layergrid-cli","rules":[]}},"results":[`)
	for i, f := range r.Findings {
		if i > 0 {
			fmt.Fprintf(&b, ",")
		}
		fmt.Fprintf(&b, `{"ruleId":%q,"level":%q,"message":{"text":%q},"locations":[{"physicalLocation":{"artifactLocation":{"uri":%q},"region":{"startLine":%d}}}]}`, f.RuleID, sarifLevel(f.Severity), f.Rationale, f.Location.Path, f.Location.Line)
	}
	fmt.Fprintf(&b, `]}]}`)
	return b.Bytes(), nil
}

func sarifLevel(sev any) string {
	switch fmt.Sprint(sev) {
	case "critical", "high":
		return "error"
	case "medium":
		return "warning"
	default:
		return "note"
	}
}
