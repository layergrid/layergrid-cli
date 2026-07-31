package report

import (
	"testing"

	json "github.com/goccy/go-json"
	"github.com/layergrid/layergrid-cli/internal/model"
	"github.com/layergrid/layergrid-cli/internal/scan"
	"github.com/layergrid/layergrid-cli/internal/trifecta"
)

func TestSARIFIsJSONAndCarriesResults(t *testing.T) {
	out, err := (SARIF{}).Format(scan.Result{Findings: []trifecta.Finding{{
		RuleID: "LG-TEST", RuleName: "Test rule", Category: "test", Severity: trifecta.SeverityHigh,
		Location:  model.Location{Path: "agent.py", Line: 7},
		Rationale: "test rationale", Fix: "test fix",
	}}})
	if err != nil {
		t.Fatal(err)
	}
	var parsed struct {
		Version string `json:"version"`
		Runs    []struct {
			Results []struct {
				RuleID string `json:"ruleId"`
			} `json:"results"`
		} `json:"runs"`
	}
	if err := json.Unmarshal(out, &parsed); err != nil {
		t.Fatal(err)
	}
	if parsed.Version != "2.1.0" || len(parsed.Runs) != 1 || len(parsed.Runs[0].Results) != 1 {
		t.Fatalf("unexpected SARIF shape: %s", string(out))
	}
}
