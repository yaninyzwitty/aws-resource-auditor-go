package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/urfave/cli/v3"
)

const (
	LAMBDA_LAST_MODIFIED string = "2006-01-02T15:04:05.000-0700"
)

func LambdaCommand() *cli.Command {
	return &cli.Command{
		Name:  "lambda",
		Usage: "Audit Lambda functions",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "old-functions",
				Usage: "find old functions not recently updated",
			},
			&cli.BoolFlag{
				Name:  "outdated-runtime",
				Usage: "find functions using outdated runtimes",
			},
		},
		Action: lambdaAction,
	}
}

func lambdaAction(ctx context.Context, cmd *cli.Command) error {
	cfg, err := ConfigFromContext(ctx)
	if err != nil {
		return fmt.Errorf("getting config: %w", err)
	}

	loader, err := AwsLoaderFromContext(ctx)
	if err != nil {
		return fmt.Errorf("getting AWS loader: %w", err)
	}

	lambdaClient, err := loader.Lambda(ctx)
	if err != nil {
		return fmt.Errorf("creating Lambda client: %w", err)
	}

	regions := []string{cfg.AWS.Region}
	if cfg.AWS.AllRegions {
		regions, err = loader.Regions(ctx)
		if err != nil {
			return fmt.Errorf("getting regions: %w", err)
		}
	}

	oldFunctions := cmd.Bool("old-functions")
	outdatedRuntime := cmd.Bool("outdated-runtime")

	if !oldFunctions && !outdatedRuntime {
		oldFunctions = true
		outdatedRuntime = true
	}

	olderThan := cfg.Thresholds.OlderThan
	if olderThan == 0 {
		olderThan = 90 * 24 * time.Hour
	}

	var results []string

	for _, region := range regions {
		if cfg.AWS.AllRegions {
			fmt.Printf("Checking region: %s\n", region)
		}

		var client *lambda.Client
		if cfg.AWS.AllRegions {
			client, err = loader.LambdaInRegion(ctx, region)
			if err != nil {
				fmt.Printf("Error creating Lambda client for region %s: %v\n", region, err)
				continue
			}
		} else {
			client = lambdaClient
		}

		if oldFunctions {
			functions, err := findOldFunctions(ctx, client, olderThan)
			if err != nil {
				fmt.Printf("Error finding old functions: %v\n", err)
			}
			results = append(results, functions...)
		}

		if outdatedRuntime {
			functions, err := findOutdatedRuntimes(ctx, client)
			if err != nil {
				fmt.Printf("Error finding outdated runtimes: %v\n", err)
			}
			results = append(results, functions...)
		}
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

func findOldFunctions(ctx context.Context, client *lambda.Client, olderThan time.Duration) ([]string, error) {
	var results []string

	paginator := lambda.NewListFunctionsPaginator(client, &lambda.ListFunctionsInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return results, fmt.Errorf("listing functions: %w", err)
		}

		for _, function := range output.Functions {
			if function.LastModified == nil {
				continue
			}

			lastModified, err := time.Parse(LAMBDA_LAST_MODIFIED, *function.LastModified)
			if err != nil {
				continue
			}

			age := time.Since(lastModified)
			if age < olderThan {
				continue
			}

			if function.FunctionName == nil {
				continue
			}

			result := fmt.Sprintf("  Function: %s - Runtime: %s - Last modified: %s ago",
				*function.FunctionName,
				string(function.Runtime),
				formatDuration(age),
			)
			results = append(results, result)
		}
	}

	return results, nil
}

func findOutdatedRuntimes(ctx context.Context, client *lambda.Client) ([]string, error) {
	var results []string

	outdatedRuntimes := map[string]bool{
		"python3.6":     true,
		"python3.7":     true,
		"nodejs10.x":    true,
		"nodejs12.x":    true,
		"nodejs14.x":    true,
		"ruby2.5":       true,
		"ruby2.7":       true,
		"java8":         true,
		"java8.al2":     true,
		"go1.x":         true,
		"dotnetcore1.0": true,
		"dotnetcore2.0": true,
		"dotnetcore2.1": true,
	}

	paginator := lambda.NewListFunctionsPaginator(client, &lambda.ListFunctionsInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return results, fmt.Errorf("listing functions: %w", err)
		}

		for _, function := range output.Functions {
			if outdatedRuntimes[string(function.Runtime)] {
				if function.FunctionName == nil {
					continue
				}
				result := fmt.Sprintf("  Function: %s - Runtime: %s (OUTDATED)",
					*function.FunctionName,
					function.Runtime,
				)
				results = append(results, result)
			}
		}
	}

	return results, nil
}
