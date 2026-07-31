package report

import (
	"testing"
	"time"

	json "github.com/goccy/go-json"
	"github.com/layergrid/layergrid-cli/internal/model"
	"github.com/layergrid/layergrid-cli/internal/scan"
	"github.com/layergrid/layergrid-cli/internal/trifecta"
	"github.com/layergrid/layergrid-cli/internal/version"
)

func TestJSONReportIncludesStableVersions(t *testing.T) {
	out, err := (JSON{}).Format(scan.Result{
		Stack:     model.Stack{Root: "/repo", ScanID: "scan", Agents: []model.Agent{}, Tools: []model.Tool{}, MCPServers: []model.MCPServer{}, Datasources: []model.Datasource{}, Models: []model.Model{}, Errors: []model.ScanError{}},
		Findings:  []trifecta.Finding{},
		Score:     trifecta.Score{Value: 100, Grade: "A", Counts: map[trifecta.Severity]int{}},
		StartedAt: time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	var parsed map[string]any
	if err := json.Unmarshal(out, &parsed); err != nil {
		t.Fatal(err)
	}
	if parsed["schemaVersion"] != version.SchemaVersion {
		t.Fatalf("schemaVersion = %v, want %s", parsed["schemaVersion"], version.SchemaVersion)
	}
	if parsed["rubricVersion"] != version.RubricVersion {
		t.Fatalf("rubricVersion = %v, want %s", parsed["rubricVersion"], version.RubricVersion)
	}
}
