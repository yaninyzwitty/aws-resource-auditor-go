package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

func LambdaCommand() *cli.Command {
	return &cli.Command{
		Name:  "lambda",
		Usage: "Audit Lambda functions",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "never-invoked",
				Usage: "find functions never invoked",
			},
			&cli.BoolFlag{
				Name:  "high-error-rate",
				Usage: "find functions with high error rates",
			},
			&cli.BoolFlag{
				Name:  "outdated-runtime",
				Usage: "find functions using outdated runtimes",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return nil
		},
	}
}
