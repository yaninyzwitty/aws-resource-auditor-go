package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

func Load(path string) (*Config, error) {

	cfg := Default()

	if path == "" {
		path = "config.yaml"
	}

	// check whether the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file %s does not exist: %w", path, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %s: %w", path, err)
	}

	return cfg, nil
}
