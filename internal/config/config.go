package config

import (
	"fmt"
	"strings"
	"time"
)

type Option func(*Config)

type Config struct {
	AWS        AWSConfig
	Output     OutputConfig
	Filter     FilterConfig
	Thresholds ThresholdConfig
	Behavior   BehaviorConfig
	Log        LogConfig
	Services   ServiceConfig
}

type AWSConfig struct {
	Profile    string `yaml:"profile"`
	Region     string `yaml:"region"`
	AllRegions bool   `yaml:"all_regions"`
	RoleARN    string `yaml:"role_arn"`
	ExternalID string `yaml:"external_id"`
}

type OutputConfig struct {
	Format  string `yaml:"format"`
	Export  string `yaml:"export"`
	NoColor bool   `yaml:"no_color"`
	Quiet   bool   `yaml:"quiet"`
}

type FilterConfig struct {
	Severity string            `yaml:"severity"`
	Tags     map[string]string `yaml:"tags"`
}

type ThresholdConfig struct {
	OlderThan time.Duration `yaml:"older_than"`
	Days      int           `yaml:"days"`
	Last      string        `yaml:"last"`

	CPUThreshold float64 `yaml:"cpu_threshold"`
	ErrorRate    float64 `yaml:"error_rate"`
}

type BehaviorConfig struct {
	Fix      bool          `yaml:"fix"`
	Watch    bool          `yaml:"watch"`
	Interval time.Duration `yaml:"interval"`
}

type LogConfig struct {
	Level   string `yaml:"level"`
	Format  string `yaml:"format"`
	Debug   bool   `yaml:"debug"`
	Verbose bool   `yaml:"verbose"`
}

type ServiceConfig struct {
	EC2 EC2Config `yaml:"ec2"`
	EBS EBSConfig `yaml:"ebs"`
	S3  S3Config  `yaml:"s3"`
	IAM IAMConfig `yaml:"iam"`

	Secrets SecretsConfig `yaml:"secrets"`
	Lambda  LambdaConfig  `yaml:"lambda"`
	RDS     RDSConfig     `yaml:"rds"`
	Costs   CostsConfig   `yaml:"costs"`
	Tags    TagsConfig    `yaml:"tags"`
}

type EC2Config struct {
	UnusedOlderThan  time.Duration `yaml:"unused_older_than"`
	OldAMIsOlderThan time.Duration `yaml:"old_amis_older_than"`
	CPUThreshold     float64       `yaml:"cpu_threshold"`
}

type EBSConfig struct {
	OldSnapshotsOlderThan time.Duration `yaml:"old_snapshots_older_than"`
}

type S3Config struct{}

type IAMConfig struct {
	StaleKeysOlderThan   time.Duration `yaml:"stale_keys_older_than"`
	UnusedRolesOlderThan time.Duration `yaml:"unused_roles_older_than"`
}

type SecretsConfig struct {
	UnusedOlderThan    time.Duration `yaml:"unused_older_than"`
	UnrotatedOlderThan time.Duration `yaml:"unrotated_older_than"`
}

type LambdaConfig struct {
	NeverInvokedOlderThan  time.Duration `yaml:"never_invoked_older_than"`
	HighErrorRateThreshold float64       `yaml:"high_error_rate_threshold"`
}

type RDSConfig struct {
	IdleDays int `yaml:"idle_days"`
}

type CostsConfig struct {
	Last string `yaml:"last"`
}

type TagsConfig struct {
	Required []string `yaml:"required"`
}

func New(opts ...Option) *Config {
	cfg := Default()

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

func WithProfile(p string) Option {
	return func(c *Config) {
		c.AWS.Profile = p
	}
}

func WithRegion(r string) Option {
	return func(c *Config) {
		c.AWS.Region = r
	}
}

func WithOutput(format string) Option {
	return func(c *Config) {
		c.Output.Format = format
	}
}

func WithOlderThan(d time.Duration) Option {
	return func(c *Config) {
		c.Thresholds.OlderThan = d
	}
}

func WithCPUThreshold(v float64) Option {
	return func(c *Config) {
		c.Thresholds.CPUThreshold = v
	}
}

func WithFix(v bool) Option {
	return func(c *Config) {
		c.Behavior.Fix = v
	}
}

func WithWatch(v bool) Option {
	return func(c *Config) {
		c.Behavior.Watch = v
	}
}

func WithInterval(d time.Duration) Option {
	return func(c *Config) {
		c.Behavior.Interval = d
	}
}

func (c *Config) Validate() error {

	if c.AWS.Region == "" && !c.AWS.AllRegions {
		return fmt.Errorf("region must be specified")
	}

	switch c.Output.Format {
	case "table", "json", "csv", "markdown":
	default:
		return fmt.Errorf("invalid output format")
	}

	return nil
}

func ParseDuration(s string) (time.Duration, error) {

	if strings.HasSuffix(s, "d") {
		var days int
		_, err := fmt.Sscanf(strings.TrimSuffix(s, "d"), "%d", &days)
		if err != nil {
			return 0, err
		}

		return time.Duration(days) * 24 * time.Hour, nil
	}

	return time.ParseDuration(s)
}
