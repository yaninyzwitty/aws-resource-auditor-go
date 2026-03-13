package config

import (
	"strings"
	"time"
)

type FlagOverrides struct {
	Profile    string
	Region     string
	AllRegions bool
	RoleARN    string
	ExternalID string

	Output  string
	Export  string
	NoColor bool
	Quiet   bool

	Severity string
	Tags     []string

	OlderThan    time.Duration
	Days         int
	Last         string
	CPUThreshold float64
	ErrorRate    float64

	Fix      bool
	Watch    bool
	Interval time.Duration

	LogLevel string
	Verbose  bool
	Debug    bool
}

func (c *Config) MergeFlags(f FlagOverrides) {

	if f.Profile != "" {
		c.AWS.Profile = f.Profile
	}

	if f.Region != "" {
		c.AWS.Region = f.Region
	}

	if f.AllRegions {
		c.AWS.AllRegions = true
	}

	if f.RoleARN != "" {
		c.AWS.RoleARN = f.RoleARN
	}

	if f.Output != "" {
		c.Output.Format = f.Output
	}

	if f.Export != "" {
		c.Output.Export = f.Export
	}

	if f.NoColor {
		c.Output.NoColor = true
	}

	if f.Quiet {
		c.Output.Quiet = true
	}

	if f.Severity != "" {
		c.Filter.Severity = f.Severity
	}

	if len(f.Tags) > 0 {
		c.Filter.Tags = make(map[string]string, len(f.Tags))
		for _, tag := range f.Tags {
			// Parse "key=value" format
			if idx := strings.Index(tag, "="); idx > 0 {
				key := tag[:idx]
				value := tag[idx+1:]
				c.Filter.Tags[key] = value
			}
		}
	}

	if f.OlderThan != 0 {
		c.Thresholds.OlderThan = Duration(f.OlderThan)
	}

	if f.Days != 0 {
		c.Thresholds.Days = f.Days
	}

	if f.Last != "" {
		c.Thresholds.Last = f.Last
	}

	if f.CPUThreshold != 0 {
		c.Thresholds.CPUThreshold = f.CPUThreshold
	}

	if f.ErrorRate != 0 {
		c.Thresholds.ErrorRate = f.ErrorRate
	}

	if f.Fix {
		c.Behavior.Fix = true
	}

	if f.Watch {
		c.Behavior.Watch = true
	}

	if f.Interval != 0 {
		c.Behavior.Interval = f.Interval
	}

	if f.LogLevel != "" {
		c.Log.Level = f.LogLevel
	}

	if f.Verbose {
		c.Log.Verbose = true
	}

	if f.Debug {
		c.Log.Debug = true
	}
}
