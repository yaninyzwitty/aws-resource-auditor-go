package config

import "time"

type GlobalConfig struct {
	// global configuration fields
	Region     string
	Profile    string
	AllRegions bool
	RoleArn    string // if any
	ExternalID string

	// Output
	Output  string // JSON | TABLE | CSV etc
	Export  string // file path to export results, if empty, print to stdout
	NoColor bool   // TODO=IMPLEMENT WITH BUBBLETEA OR SOMETHING disable color output
	Quiet   bool   // disable all output except errors

	// Filtering
	Severity string            // critical|high|medium|low
	Tags     map[string]string // --tags Name=witty,Env=prod

	// shared flags
	OlderThan    time.Duration // --older-than 90d
	Days         int           // --days 90 (alternative to older-than, parsed as hours internally)
	Last         string
	CPUThreshold float64 // --cpu-threshold 80.0
	Threshold    float64 // generic % threshold (error rate, utilization, etc.)

	// Behavior
	Fix      bool          // dry-run remediation
	Watch    bool          // daemon mode
	Interval time.Duration // --interval 6h

	// Debug
	Verbose  bool
	LogLevel string
	Debug    bool
}
