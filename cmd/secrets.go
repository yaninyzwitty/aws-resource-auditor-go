package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/urfave/cli/v3"
)

func SecretsCommand() *cli.Command {
	return &cli.Command{
		Name:  "secrets",
		Usage: "Audit Secrets Manager secrets",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "unrotated",
				Usage: "find unrotated secrets",
			},
		},
		Action: secretsAction,
	}
}

func secretsAction(ctx context.Context, cmd *cli.Command) error {
	loader, err := AwsLoaderFromContext(ctx)
	if err != nil {
		return fmt.Errorf("getting AWS loader: %w", err)
	}

	secretsClient, err := loader.SecretsManager(ctx)
	if err != nil {
		return fmt.Errorf("creating Secrets Manager client: %w", err)
	}

	unrotated := cmd.Bool("unrotated")
	// default to all checks when no flags specified
	if !unrotated {
		unrotated = true
	}

	cfg, _ := ConfigFromContext(ctx)
	olderThan := cfg.Services.Secrets.UnrotatedOlderThan
	if olderThan == 0 {
		olderThan = 90 * 24 * time.Hour
	}

	var results []string

	if unrotated {
		secrets, err := findUnrotatedSecrets(ctx, secretsClient, olderThan)
		if err != nil {
			fmt.Printf("Error finding unrotated secrets: %v\n", err)
		}
		results = append(results, secrets...)
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

func findUnrotatedSecrets(ctx context.Context, client *secretsmanager.Client, olderThan time.Duration) ([]string, error) {
	var results []string

	paginator := secretsmanager.NewListSecretsPaginator(client, &secretsmanager.ListSecretsInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return results, fmt.Errorf("listing secrets: %w", err)
		}

		for _, secret := range output.SecretList {
			if secret.LastRotatedDate == nil {
				if secret.CreatedDate == nil {
					continue
				}
				age := time.Since(*secret.CreatedDate)
				if age >= olderThan {
					if secret.Name == nil {
						continue
					}
					result := fmt.Sprintf("  Secret: %s - NEVER ROTATED - Created: %s ago",
						*secret.Name,
						formatDuration(age),
					)
					results = append(results, result)
				}
				continue
			}

			age := time.Since(*secret.LastRotatedDate)
			if age >= olderThan {
				result := fmt.Sprintf("  Secret: %s - Last rotated: %s ago",
					*secret.Name,
					formatDuration(age),
				)
				results = append(results, result)
			}
		}
	}

	return results, nil
}
