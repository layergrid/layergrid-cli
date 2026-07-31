package trifecta

import "testing"

func TestComputeCapsRepeatedRuleImpact(t *testing.T) {
	findings := make([]Finding, 6)
	for i := range findings {
		findings[i] = Finding{RuleID: "R", Severity: SeverityHigh, ScoreImpact: -15}
	}
	score := Compute(findings)
	if score.Value != 25 {
		t.Fatalf("score = %d, want 25", score.Value)
	}
	if score.Counts[SeverityHigh] != 6 {
		t.Fatalf("high count = %d, want 6", score.Counts[SeverityHigh])
	}
}
