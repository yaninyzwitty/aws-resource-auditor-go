package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

func RDSCommand() *cli.Command {
	return &cli.Command{
		Name:  "rds",
		Usage: "Audit RDS instances",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "idle",
				Usage: "find idle instances",
			},
			&cli.BoolFlag{
				Name:  "unencrypted",
				Usage: "find unencrypted instances",
			},
			&cli.BoolFlag{
				Name:  "public",
				Usage: "find publicly accessible instances",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return nil
		},
	}
}
