package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

func S3Command() *cli.Command {
	return &cli.Command{
		Name:  "s3",
		Usage: "Audit S3 buckets",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "public",
				Usage: "find public buckets",
			},
			&cli.BoolFlag{
				Name:  "unencrypted",
				Usage: "find unencrypted buckets",
			},
			&cli.BoolFlag{
				Name:  "versioning",
				Usage: "check versioning status",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return nil
		},
	}
}
