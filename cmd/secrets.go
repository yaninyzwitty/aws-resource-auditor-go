package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

func SecretsCommand() *cli.Command {
	return &cli.Command{
		Name:  "secrets",
		Usage: "Audit Secrets Manager secrets",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "unused",
				Usage: "find unused secrets",
			},
			&cli.BoolFlag{
				Name:  "unrotated",
				Usage: "find unrotated secrets",
			},
			&cli.BoolFlag{
				Name:  "public",
				Usage: "find publicly accessible secrets",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return nil
		},
	}
}
