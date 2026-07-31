package trifecta

import (
	"io/fs"
	"sort"

	builtinrules "github.com/layergrid/layergrid-cli/rules"
	"gopkg.in/yaml.v3"
)

func LoadBuiltinRules() ([]Rule, error) {
	var rules []Rule
	err := fs.WalkDir(builtinrules.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || len(path) < 6 || path[len(path)-5:] != ".yaml" {
			return nil
		}
		data, err := builtinrules.FS.ReadFile(path)
		if err != nil {
			return err
		}
		var rule Rule
		if err := yaml.Unmarshal(data, &rule); err != nil {
			return err
		}
		rules = append(rules, rule)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(rules, func(i, j int) bool { return rules[i].ID < rules[j].ID })
	return rules, nil
}
