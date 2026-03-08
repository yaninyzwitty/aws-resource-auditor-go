package cmd

import (
	"time"

	"github.com/urfave/cli/v3"
)

func GlobalFlags() []cli.Flag {
	return []cli.Flag{

		// config file
		&cli.StringFlag{
			Name:  "config",
			Usage: "config file path",
		},

		// AWS
		&cli.StringFlag{
			Name:        "profile",
			Aliases:     []string{"p"},
			Usage:       "AWS named profile",
			Destination: &globalConfig.Profile,
			Sources:     cli.EnvVars("AWS_PROFILE", "PROFILE"),
		},

		&cli.StringFlag{
			Name:        "region",
			Aliases:     []string{"r"},
			Usage:       "AWS region",
			Destination: &globalConfig.Region,
			Sources:     cli.EnvVars("AWS_REGION", "REGION"),
		},

		&cli.BoolFlag{
			Name:        "all-regions",
			Usage:       "scan all regions",
			Destination: &globalConfig.AllRegions,
		},

		&cli.StringFlag{
			Name:        "role-arn",
			Usage:       "IAM role ARN to assume",
			Sources:     cli.EnvVars("AUDIT_ROLE_ARN"),
			Destination: &globalConfig.RoleARN,
		},

		&cli.StringFlag{
			Name:        "external-id",
			Usage:       "external ID for role assumption",
			Destination: &globalConfig.ExternalID,
		},

		// OUTPUT
		&cli.StringFlag{
			Name:        "output",
			Aliases:     []string{"o"},
			Value:       "table",
			Usage:       "output format: table|json|csv|markdown",
			Sources:     cli.EnvVars("AUDIT_OUTPUT"),
			Destination: &globalConfig.Output,
		},

		&cli.StringFlag{
			Name:        "export",
			Usage:       "write results to file",
			Sources:     cli.EnvVars("AUDIT_EXPORT"),
			Destination: &globalConfig.Export,
		},

		&cli.BoolFlag{
			Name:        "no-color",
			Usage:       "disable ANSI color output",
			Sources:     cli.EnvVars("NO_COLOR"),
			Destination: &globalConfig.NoColor,
		},

		&cli.BoolFlag{
			Name:        "quiet",
			Aliases:     []string{"q"},
			Destination: &globalConfig.Quiet,
		},

		// FILTERING
		&cli.StringFlag{
			Name:        "severity",
			Value:       "low",
			Usage:       "minimum severity: critical|high|medium|low",
			Sources:     cli.EnvVars("AUDIT_SEVERITY"),
			Destination: &globalConfig.Severity,
		},

		// THRESHOLDS
		&cli.DurationFlag{
			Name:        "older-than",
			Usage:       "age threshold e.g. 90d",
			Destination: &globalConfig.OlderThan,
		},

		&cli.IntFlag{
			Name:        "days",
			Usage:       "lookback window in days",
			Destination: &globalConfig.Days,
		},

		&cli.StringFlag{
			Name:        "last",
			Usage:       "time window e.g. 30d",
			Destination: &globalConfig.Last,
		},

		&cli.Float64Flag{
			Name:        "cpu-threshold",
			Value:       5,
			Usage:       "CPU % threshold",
			Destination: &globalConfig.CPUThreshold,
		},

		&cli.Float64Flag{
			Name:        "threshold",
			Usage:       "generic utilization threshold",
			Destination: &globalConfig.ErrorRate,
		},

		// BEHAVIOR
		&cli.BoolFlag{
			Name:        "fix",
			Usage:       "emit remediation suggestions",
			Destination: &globalConfig.Fix,
		},

		&cli.BoolFlag{
			Name:        "watch",
			Usage:       "daemon mode",
			Destination: &globalConfig.Watch,
		},

		&cli.DurationFlag{
			Name:        "interval",
			Value:       6 * time.Hour,
			Usage:       "watch interval",
			Destination: &globalConfig.Interval,
		},

		// LOGGING
		&cli.BoolFlag{
			Name:        "verbose",
			Aliases:     []string{"v"},
			Destination: &globalConfig.Verbose,
		},

		&cli.BoolFlag{
			Name:        "debug",
			Destination: &globalConfig.Debug,
		},

		&cli.StringFlag{
			Name:        "log-level",
			Value:       "info",
			Sources:     cli.EnvVars("AUDIT_LOG_LEVEL"),
			Destination: &globalConfig.LogLevel,
		},

		// TAG FILTER
		&cli.StringSliceFlag{
			Name:  "tags",
			Usage: "filter tags e.g. --tags Name=witty",
		},
	}
}
