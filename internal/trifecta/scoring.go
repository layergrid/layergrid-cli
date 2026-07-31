package trifecta

func Compute(findings []Finding) Score {
	total := 100
	countedByRule := map[string]int{}
	counts := map[Severity]int{
		SeverityCritical: 0,
		SeverityHigh:     0,
		SeverityMedium:   0,
		SeverityLow:      0,
	}
	for _, f := range findings {
		counts[f.Severity]++
		if f.Confidence == "low" {
			continue
		}
		if countedByRule[f.RuleID] < 5 {
			total += f.ScoreImpact
			countedByRule[f.RuleID]++
		}
	}
	if total < 0 {
		total = 0
	}
	return Score{Value: total, Grade: gradeFor(total), Counts: counts}
}

func gradeFor(v int) string {
	switch {
	case v >= 90:
		return "A"
	case v >= 75:
		return "B"
	case v >= 60:
		return "C"
	case v >= 40:
		return "D"
	default:
		return "F"
	}
}
