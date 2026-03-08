package config

import (
	"os"

	"go.yaml.in/yaml/v3"
)

func Load(path string) (*Config, error) {

	cfg := Default()

	if path == "" {
		path = "config.yaml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, nil
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
