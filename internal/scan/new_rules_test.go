package scan

import (
	"path/filepath"
	"testing"
)

func TestNewRuleFixtures(t *testing.T) {
	cases := []struct {
		ruleID string
	}{
		{ruleID: "LG-MCP-CREDENTIAL-IN-ENV"},
		{ruleID: "LG-TOOL-UNICODE-HIDDEN"},
		{ruleID: "LG-MCP-EGRESS-UNBOUND"},
		{ruleID: "LG-AGENT-INBOX-EXFIL"},
		{ruleID: "LG-MEMORY-EXTERNAL-WRITE"},
	}
	for _, tc := range cases {
		t.Run(tc.ruleID+"/vuln", func(t *testing.T) {
			result, err := Run(filepath.Join("..", "..", "testdata", tc.ruleID, "vuln"), Options{})
			if err != nil {
				t.Fatal(err)
			}
			if !hasFinding(result, tc.ruleID) {
				t.Fatalf("expected finding %s; got %v", tc.ruleID, findingIDs(result))
			}
		})
		t.Run(tc.ruleID+"/safe", func(t *testing.T) {
			result, err := Run(filepath.Join("..", "..", "testdata", tc.ruleID, "safe"), Options{})
			if err != nil {
				t.Fatal(err)
			}
			if hasFinding(result, tc.ruleID) {
				t.Fatalf("did not expect finding %s; got %v", tc.ruleID, findingIDs(result))
			}
		})
	}
}

func hasFinding(result Result, ruleID string) bool {
	for _, finding := range result.Findings {
		if finding.RuleID == ruleID {
			return true
		}
	}
	return false
}

func findingIDs(result Result) []string {
	ids := make([]string, 0, len(result.Findings))
	for _, finding := range result.Findings {
		ids = append(ids, finding.RuleID)
	}
	return ids
}
