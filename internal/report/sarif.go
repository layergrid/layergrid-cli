package report

import (
	"fmt"

	json "github.com/goccy/go-json"
	"github.com/layergrid/layergrid-cli/internal/scan"
)

type SARIF struct{}

type sarifReport struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name           string      `json:"name"`
	InformationURI string      `json:"informationUri"`
	Rules          []sarifRule `json:"rules"`
}

type sarifRule struct {
	ID               string              `json:"id"`
	Name             string              `json:"name"`
	ShortDescription sarifText           `json:"shortDescription"`
	FullDescription  sarifText           `json:"fullDescription"`
	HelpURI          string              `json:"helpUri,omitempty"`
	Properties       sarifRuleProperties `json:"properties"`
}

type sarifRuleProperties struct {
	Severity string `json:"severity"`
	Category string `json:"category"`
}

type sarifResult struct {
	RuleID    string          `json:"ruleId"`
	RuleIndex int             `json:"ruleIndex,omitempty"`
	Level     string          `json:"level"`
	Message   sarifText       `json:"message"`
	Locations []sarifLocation `json:"locations"`
}

type sarifText struct {
	Text string `json:"text"`
}

type sarifLocation struct {
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}

type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifactLocation `json:"artifactLocation"`
	Region           sarifRegion           `json:"region"`
}

type sarifArtifactLocation struct {
	URI string `json:"uri"`
}

type sarifRegion struct {
	StartLine int `json:"startLine,omitempty"`
}

func (SARIF) Format(r scan.Result) ([]byte, error) {
	ruleIndex := map[string]int{}
	rules := make([]sarifRule, 0)
	results := make([]sarifResult, 0, len(r.Findings))
	for _, f := range r.Findings {
		idx, ok := ruleIndex[f.RuleID]
		if !ok {
			idx = len(rules)
			ruleIndex[f.RuleID] = idx
			helpURI := ""
			if len(f.References) > 0 {
				helpURI = f.References[0]
			}
			rules = append(rules, sarifRule{
				ID: f.RuleID, Name: f.RuleName,
				ShortDescription: sarifText{Text: f.RuleName},
				FullDescription:  sarifText{Text: f.Rationale},
				HelpURI:          helpURI,
				Properties:       sarifRuleProperties{Severity: string(f.Severity), Category: f.Category},
			})
		}
		results = append(results, sarifResult{
			RuleID: f.RuleID, RuleIndex: idx, Level: sarifLevel(f.Severity),
			Message: sarifText{Text: fmt.Sprintf("%s Fix: %s", f.Rationale, f.Fix)},
			Locations: []sarifLocation{{
				PhysicalLocation: sarifPhysicalLocation{
					ArtifactLocation: sarifArtifactLocation{URI: f.Location.Path},
					Region:           sarifRegion{StartLine: f.Location.Line},
				},
			}},
		})
	}
	return json.MarshalIndent(sarifReport{
		Version: "2.1.0",
		Schema:  "https://json.schemastore.org/sarif-2.1.0.json",
		Runs: []sarifRun{{
			Tool: sarifTool{Driver: sarifDriver{
				Name: "LayerGrid", InformationURI: "https://github.com/layergrid/layergrid-cli", Rules: rules,
			}},
			Results: results,
		}},
	}, "", "  ")
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
