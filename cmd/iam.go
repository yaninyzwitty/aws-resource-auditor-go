package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iam"
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
		Action: iamAction,
	}
}

func iamAction(ctx context.Context, cmd *cli.Command) error {
	loader, err := AwsLoaderFromContext(ctx)
	if err != nil {
		return fmt.Errorf("getting AWS loader: %w", err)
	}

	iamClient, err := loader.IAM(ctx)
	if err != nil {
		return fmt.Errorf("creating IAM client: %w", err)
	}

	staleKeys := cmd.Bool("stale-keys")
	unusedRoles := cmd.Bool("unused-roles")
	passwordPolicy := cmd.Bool("password-policy")

	if !staleKeys && !unusedRoles && !passwordPolicy {
		staleKeys = true
		unusedRoles = true
		passwordPolicy = true
	}

	var results []string

	if staleKeys {
		keys, err := findStaleAccessKeys(ctx, iamClient, 90*24*time.Hour)
		if err != nil {
			fmt.Printf("Error finding stale keys: %v\n", err)
		}
		results = append(results, keys...)
	}

	if unusedRoles {
		roles, err := findUnusedRoles(ctx, iamClient, 90*24*time.Hour)
		if err != nil {
			fmt.Printf("Error finding unused roles: %v\n", err)
		}
		results = append(results, roles...)
	}

	if passwordPolicy {
		policy, err := checkPasswordPolicy(ctx, iamClient)
		if err != nil {
			fmt.Printf("Error checking password policy: %v\n", err)
		} else if policy != "" {
			results = append(results, policy)
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

func findStaleAccessKeys(ctx context.Context, client *iam.Client, olderThan time.Duration) ([]string, error) {
	var results []string

	paginator := iam.NewListUsersPaginator(client, &iam.ListUsersInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return results, fmt.Errorf("listing users: %w", err)
		}

		for _, user := range output.Users {
			keysOutput, err := client.ListAccessKeys(ctx, &iam.ListAccessKeysInput{
				UserName: user.UserName,
			})
			if err != nil {
				continue
			}

			for _, key := range keysOutput.AccessKeyMetadata {
				if key.CreateDate == nil {
					continue
				}
				age := time.Since(*key.CreateDate)

				if age < olderThan {
					continue
				}

				result := fmt.Sprintf("  User: %s - Key: %s - Status: %s - Created: %s ago",
					*user.UserName,
					*key.AccessKeyId,
					string(key.Status),
					formatDuration(age),
				)
				results = append(results, result)
			}
		}
	}

	return results, nil
}

func findUnusedRoles(ctx context.Context, client *iam.Client, olderThan time.Duration) ([]string, error) {
	var results []string

	paginator := iam.NewListRolesPaginator(client, &iam.ListRolesInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return results, fmt.Errorf("listing roles: %w", err)
		}

		for _, role := range output.Roles {
			if role.RoleLastUsed == nil || role.RoleLastUsed.LastUsedDate == nil {
				if role.CreateDate == nil {
					continue
				}
				age := time.Since(*role.CreateDate)
				if age < olderThan {
					continue
				}
				result := fmt.Sprintf("  Role: %s - NEVER USED - Created: %s ago",
					*role.RoleName,
					formatDuration(age),
				)
				results = append(results, result)
				continue
			}

			age := time.Since(*role.RoleLastUsed.LastUsedDate)
			if age < olderThan {
				continue
			}

			result := fmt.Sprintf("  Role: %s - Last used: %s ago",
				*role.RoleName,
				formatDuration(age),
			)
			results = append(results, result)
		}
	}

	return results, nil
}

func checkPasswordPolicy(ctx context.Context, client *iam.Client) (string, error) {
	output, err := client.GetAccountPasswordPolicy(ctx, &iam.GetAccountPasswordPolicyInput{})
	if err != nil {
		if isNoSuchEntity(err) {
			return "  Password Policy Issues:\n    - No custom password policy configured (using AWS defaults)", nil
		}
		return "", err
	}

	policy := output.PasswordPolicy
	issues := []string{}

	if policy.MinimumPasswordLength != nil && *policy.MinimumPasswordLength < 14 {
		issues = append(issues, fmt.Sprintf("Min length: %d (recommended: 14)", *policy.MinimumPasswordLength))
	}

	if !policy.RequireUppercaseCharacters {
		issues = append(issues, "Missing uppercase requirement")
	}

	if !policy.RequireLowercaseCharacters {
		issues = append(issues, "Missing lowercase requirement")
	}

	if !policy.RequireNumbers {
		issues = append(issues, "Missing number requirement")
	}

	if !policy.RequireSymbols {
		issues = append(issues, "Missing symbol requirement")
	}

	if policy.MaxPasswordAge == nil || *policy.MaxPasswordAge > 90 {
		issues = append(issues, "Password expiry not set or too long")
	}

	if len(issues) == 0 {
		return "", nil
	}

	return "  Password Policy Issues:\n    - " + strings.Join(issues, "\n    - "), nil
}

func isNoSuchEntity(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "NoSuchEntity")
}
