package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/yaninyzwitty/aws-resource-auditor-go/internal/aws"
	"github.com/yaninyzwitty/aws-resource-auditor-go/internal/config"
)

func NewCliCommand() *cli.Command {
	return &cli.Command{
		Name:  "aws-resource-auditor",
		Usage: "Audit AWS resources",

		Flags: GlobalFlags(),

		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			cfg, err := config.Load(cmd.String("config"))
			if err != nil {
				return ctx, fmt.Errorf("config load: %w", err)
			}

			cfg.MergeFlags(globalConfig)

			if err := cfg.Validate(); err != nil {
				return ctx, fmt.Errorf("config validation: %w", err)
			}

			loader, err := aws.NewLoader(cfg.AWS)
			if err != nil {
				return ctx, fmt.Errorf("AWS loader: %w", err)
			}

			ctx = context.WithValue(ctx, ConfigKey, cfg)
			ctx = context.WithValue(ctx, AwsLoaderKey, loader)

			return ctx, nil
		},

		Commands: []*cli.Command{
			EC2Command(),
			S3Command(),
			IAMCommand(),
			LambdaCommand(),
			RDSCommand(),
			SecretsCommand(),
		},
	}
}
