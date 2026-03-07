package cmd

import (
	"github.com/urfave/cli/v3"
)

func NewCliCommand() *cli.Command {
	return &cli.Command{
		Name:        "audit",
		Description: "An audit system for auditing your AWS resources",
		Usage:       "Audit your AWS resources",

		// global flags
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "region",
				Aliases: []string{"r"},
				Usage:   "Aws target region",
				Sources: cli.EnvVars("AWS_REGION"),
				Value:   "us-east-1",
			},
			&cli.StringFlag{
				Name:    "profile",
				Aliases: []string{"p"},
				Usage:   "Aws named profile",
				Sources: cli.EnvVars("AWS_PROFILE"),
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output format: table, JSON, CSV, MarkDown",
				Value:   "table",
			},
			&cli.BoolFlag{
				Name:    "all-regions",
				Usage:   "Scan all regions. If set, the --region flag will be ignored.",
				Value:   false,
				Sources: cli.EnvVars("ALL_REGIONS"),
			},
		},
		Commands: []*cli.Command{
			// commands might live here
		},
	}
}
