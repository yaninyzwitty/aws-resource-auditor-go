package config

import "time"

func Default() *Config {
	return &Config{
		AWS: AWSConfig{
			Region: "us-east-1",
		},

		Output: OutputConfig{
			Format: "table",
		},

		Filter: FilterConfig{
			Tags: map[string]string{},
		},

		Thresholds: ThresholdConfig{
			OlderThan:    90 * 24 * time.Hour,
			CPUThreshold: 20,
		},

		Behavior: BehaviorConfig{
			Interval: 6 * time.Hour,
		},

		Log: LogConfig{
			Level:  "info",
			Format: "text",
		},
	}
}
