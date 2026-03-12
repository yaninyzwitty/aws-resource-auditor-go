package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/rds"
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
		Action: rdsAction,
	}
}

func rdsAction(ctx context.Context, cmd *cli.Command) error {
	loader, err := AwsLoaderFromContext(ctx)
	if err != nil {
		return fmt.Errorf("getting AWS loader: %w", err)
	}

	rdsClient, err := loader.RDS(ctx)
	if err != nil {
		return fmt.Errorf("creating RDS client: %w", err)
	}

	idle := cmd.Bool("idle")
	unencrypted := cmd.Bool("unencrypted")
	public := cmd.Bool("public")

	if !idle && !unencrypted && !public {
		idle = true
		unencrypted = true
		public = true
	}

	cfg, _ := ConfigFromContext(ctx)
	idleDays := cfg.Services.RDS.IdleDays
	if idleDays == 0 {
		idleDays = 30
	}

	var results []string

	if idle {
		instances, err := findIdleInstances(ctx, rdsClient, idleDays)
		if err != nil {
			fmt.Printf("Error finding idle instances: %v\n", err)
		}
		results = append(results, instances...)
	}

	if unencrypted {
		instances, err := findUnencryptedRDS(ctx, rdsClient)
		if err != nil {
			fmt.Printf("Error finding unencrypted instances: %v\n", err)
		}
		results = append(results, instances...)
	}

	if public {
		instances, err := findPublicRDS(ctx, rdsClient)
		if err != nil {
			fmt.Printf("Error finding public instances: %v\n", err)
		}
		results = append(results, instances...)
	}

	if len(results) == 0 {
		fmt.Println("No issues found")
		return nil
	}

	fmt.Println("Findings:")
	for _, r := range results {
		fmt.Println(r)
	}

	return nil
}

func findIdleInstances(ctx context.Context, client *rds.Client, idleDays int) ([]string, error) {
	var results []string

	paginator := rds.NewDescribeDBInstancesPaginator(client, &rds.DescribeDBInstancesInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return results, fmt.Errorf("describing instances: %w", err)
		}

		for _, instance := range output.DBInstances {
			if instance.DBInstanceStatus == nil || *instance.DBInstanceStatus != "available" {
				continue
			}

			createTime := time.Now()
			if instance.InstanceCreateTime != nil {
				createTime = *instance.InstanceCreateTime
			}

			daysIdle := int(time.Since(createTime).Hours() / 24)
			if daysIdle < idleDays {
				continue
			}

			result := fmt.Sprintf("  Instance: %s (%s) - Engine: %s - Age: %d days",
				*instance.DBInstanceIdentifier,
				*instance.DBInstanceClass,
				*instance.Engine,
				daysIdle,
			)
			results = append(results, result)
		}
	}

	return results, nil
}

func findUnencryptedRDS(ctx context.Context, client *rds.Client) ([]string, error) {
	var results []string

	paginator := rds.NewDescribeDBInstancesPaginator(client, &rds.DescribeDBInstancesInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return results, fmt.Errorf("describing instances: %w", err)
		}

		for _, instance := range output.DBInstances {
			if instance.StorageEncrypted != nil && *instance.StorageEncrypted {
				continue
			}

			result := fmt.Sprintf("  Instance: %s (%s) - UNENCRYPTED",
				*instance.DBInstanceIdentifier,
				*instance.Engine,
			)
			results = append(results, result)
		}
	}

	return results, nil
}

func findPublicRDS(ctx context.Context, client *rds.Client) ([]string, error) {
	var results []string

	paginator := rds.NewDescribeDBInstancesPaginator(client, &rds.DescribeDBInstancesInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return results, fmt.Errorf("describing instances: %w", err)
		}

		for _, instance := range output.DBInstances {
			if instance.PubliclyAccessible == nil || !*instance.PubliclyAccessible {
				continue
			}

			result := fmt.Sprintf("  Instance: %s (%s) - PUBLIC",
				*instance.DBInstanceIdentifier,
				*instance.Engine,
			)
			results = append(results, result)
		}
	}

	return results, nil
}
