package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/yaninyzwitty/aws-resource-auditor-go/internal/config"
)

var (
	globalConfig config.FlagOverrides
	Cfg          *config.Config
)

func NewCliCommand() *cli.Command {
	return &cli.Command{
		Name:  "aws-resource-auditor",
		Usage: "Audit AWS resources",

		Flags: GlobalFlags(),

		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {

			// load config file
			cfg, err := config.Load(cmd.String("config"))
			if err != nil {
				return ctx, fmt.Errorf("config load: %w", err)
			}

			// merge CLI flags
			cfg.MergeFlags(globalConfig)

			// validate config
			if err := cfg.Validate(); err != nil {
				return ctx, fmt.Errorf("config validation: %w", err)
			}

			Cfg = cfg

			return ctx, nil
		},

		Commands: []*cli.Command{},
	}
}
