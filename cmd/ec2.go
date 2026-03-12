package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

func EC2Command() *cli.Command {
	return &cli.Command{
		Name:  "ec2",
		Usage: "Audit EC2 instances and resources",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "unused",
				Usage: "find unused instances",
			},
			&cli.BoolFlag{
				Name:  "old-amis",
				Usage: "find old unused AMIs",
			},
			&cli.BoolFlag{
				Name:  "unencrypted",
				Usage: "find unencrypted volumes",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return nil
		},
	}
}
