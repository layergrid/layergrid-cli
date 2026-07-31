package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version    int      `yaml:"version"`
	Include    []string `yaml:"include"`
	Exclude    []string `yaml:"exclude"`
	Frameworks []string `yaml:"frameworks"`
	Rules      Rules    `yaml:"rules"`
	FailOn     string   `yaml:"fail_on"`
}

type Rules struct {
	Disable    []string `yaml:"disable"`
	Categories []string `yaml:"categories"`
}

func Load(root, explicitPath string) (Config, bool, error) {
	path := explicitPath
	if path == "" {
		path = filepath.Join(root, ".layergrid.yaml")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) && explicitPath == "" {
			return Config{}, false, nil
		}
		return Config{}, false, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, false, err
	}
	return cfg, true, nil
}
