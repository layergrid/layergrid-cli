package report

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/layergrid/layergrid/internal/scan"
)

type Markdown struct{}

func (Markdown) Format(r scan.Result) ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "# LayerGrid Scan\n\n")
	fmt.Fprintf(&b, "**Score:** %d / 100 (%s)\n\n", r.Score.Value, r.Score.Grade)
	fmt.Fprintf(&b, "| Severity | Rule | Location | Fix |\n")
	fmt.Fprintf(&b, "|---|---|---|---|\n")
	for _, f := range r.Findings {
		fmt.Fprintf(&b, "| %s | `%s` | `%s:%d` | %s |\n", strings.ToUpper(string(f.Severity)), f.RuleID, f.Location.Path, f.Location.Line, strings.ReplaceAll(f.Fix, "|", "\\|"))
	}
	return b.Bytes(), nil
}
