package trifecta

import "github.com/layergrid/layergrid-cli/internal/model"

type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
)

type Rule struct {
	ID          string   `yaml:"id" json:"id"`
	Name        string   `yaml:"name" json:"name"`
	Severity    Severity `yaml:"severity" json:"severity"`
	Category    string   `yaml:"category" json:"category"`
	Description string   `yaml:"description" json:"description"`
	References  []string `yaml:"references" json:"references"`
	Fix         string   `yaml:"fix" json:"fix"`
	ScoreImpact int      `yaml:"score_impact" json:"scoreImpact"`
}

type Subject struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type PathNode struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type Finding struct {
	ID          string         `json:"id"`
	RuleID      string         `json:"ruleId"`
	RuleName    string         `json:"ruleName"`
	Category    string         `json:"category"`
	Severity    Severity       `json:"severity"`
	Subject     Subject        `json:"subject"`
	Path        []PathNode     `json:"path,omitempty"`
	Location    model.Location `json:"location"`
	Fix         string         `json:"fix"`
	References  []string       `json:"references,omitempty"`
	ScoreImpact int            `json:"scoreImpact"`
	Rationale   string         `json:"rationale"`
}

type Score struct {
	Value  int              `json:"value"`
	Grade  string           `json:"grade"`
	Counts map[Severity]int `json:"counts"`
}
