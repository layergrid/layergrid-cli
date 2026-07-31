package report

import (
	json "github.com/goccy/go-json"
	"github.com/layergrid/layergrid-cli/internal/scan"
	"github.com/layergrid/layergrid-cli/internal/version"
)

type JSON struct{}

type jsonReport struct {
	SchemaVersion string   `json:"schemaVersion"`
	RubricVersion string   `json:"rubricVersion"`
	ToolVersion   string   `json:"toolVersion"`
	Scan          scanMeta `json:"scan"`
	Stack         any      `json:"stack"`
	Findings      any      `json:"findings"`
	Score         any      `json:"score"`
}

type scanMeta struct {
	ID         string `json:"id"`
	Root       string `json:"root"`
	StartedAt  string `json:"startedAt"`
	DurationMs int64  `json:"durationMs"`
}

func (JSON) Format(r scan.Result) ([]byte, error) {
	return json.MarshalIndent(jsonReport{
		SchemaVersion: version.SchemaVersion,
		RubricVersion: version.RubricVersion,
		ToolVersion:   version.Version,
		Scan: scanMeta{
			ID: r.Stack.ScanID, Root: r.Stack.Root,
			StartedAt:  r.StartedAt.Format("2006-01-02T15:04:05Z07:00"),
			DurationMs: r.Duration.Milliseconds(),
		},
		Stack: r.Stack, Findings: r.Findings, Score: r.Score,
	}, "", "  ")
}
