package scan

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/layergrid/layergrid/internal/detectors"
	"github.com/layergrid/layergrid/internal/graph"
	"github.com/layergrid/layergrid/internal/model"
	"github.com/layergrid/layergrid/internal/trifecta"
)

type Result struct {
	Stack     model.Stack        `json:"stack"`
	Graph     *graph.Graph       `json:"graph,omitempty"`
	Findings  []trifecta.Finding `json:"findings"`
	Score     trifecta.Score     `json:"score"`
	Duration  time.Duration      `json:"duration"`
	StartedAt time.Time          `json:"startedAt"`
}

func Run(path string, opts Options) (Result, error) {
	if path == "" {
		path = "."
	}
	root, err := filepath.Abs(path)
	if err != nil {
		return Result{}, err
	}
	info, err := os.Stat(root)
	if err != nil {
		return Result{}, err
	}
	if !info.IsDir() {
		return Result{}, fmt.Errorf("%s is not a directory", root)
	}

	start := time.Now().UTC()
	stack := model.Stack{
		Root:        root,
		ScanID:      model.StableID("scan", root, start.Format(time.RFC3339Nano)),
		ScannedAt:   start,
		Agents:      []model.Agent{},
		Tools:       []model.Tool{},
		MCPServers:  []model.MCPServer{},
		Datasources: []model.Datasource{},
		Models:      []model.Model{},
		Errors:      []model.ScanError{},
	}
	for _, detector := range detectors.Registry() {
		if !enabledDetector(detector, opts.Frameworks) {
			continue
		}
		runDetector(root, detector, &stack)
	}
	g := graph.Build(&stack)
	rules, err := trifecta.LoadBuiltinRules()
	if err != nil {
		return Result{}, err
	}
	rules = filterRules(rules, opts.Rules)
	findings := trifecta.Engine{Rules: rules}.Evaluate(&stack, g)
	score := trifecta.Compute(findings)
	return Result{Stack: stack, Graph: g, Findings: findings, Score: score, Duration: time.Since(start), StartedAt: start}, nil
}

func runDetector(root string, detector detectors.Detector, stack *model.Stack) {
	defer func() {
		if r := recover(); r != nil {
			stack.Errors = append(stack.Errors, model.ScanError{
				Detector: detector.Name(),
				Message:  fmt.Sprintf("detector recovered from panic: %v", r),
			})
		}
	}()
	if err := detector.Detect(root, stack); err != nil {
		stack.Errors = append(stack.Errors, model.ScanError{Detector: detector.Name(), Message: err.Error()})
	}
}

func enabledDetector(detector detectors.Detector, filters []string) bool {
	if len(filters) == 0 {
		return true
	}
	for _, f := range filters {
		if f == detector.Name() || f == string(detector.Framework()) {
			return true
		}
	}
	return false
}

func filterRules(rules []trifecta.Rule, filters []string) []trifecta.Rule {
	if len(filters) == 0 {
		return rules
	}
	var out []trifecta.Rule
	for _, rule := range rules {
		for _, f := range filters {
			if f == rule.ID || f == rule.Category {
				out = append(out, rule)
				break
			}
		}
	}
	return out
}

func ShouldFail(score trifecta.Score, failOn string) bool {
	if failOn == "" || failOn == "never" {
		return false
	}
	ranks := map[string]int{"critical": 0, "high": 1, "medium": 2, "low": 3}
	threshold, ok := ranks[failOn]
	if !ok {
		return false
	}
	for sev, count := range score.Counts {
		if count > 0 && ranks[string(sev)] <= threshold {
			return true
		}
	}
	return false
}
