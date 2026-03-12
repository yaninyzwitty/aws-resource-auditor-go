package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

func IAMCommand() *cli.Command {
	return &cli.Command{
		Name:  "iam",
		Usage: "Audit IAM users, roles, and policies",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "stale-keys",
				Usage: "find stale access keys",
			},
			&cli.BoolFlag{
				Name:  "unused-roles",
				Usage: "find unused IAM roles",
			},
			&cli.BoolFlag{
				Name:  "password-policy",
				Usage: "check password policy",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return nil
		},
	}
}
