package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
		Action: s3Action,
	}
}

func s3Action(ctx context.Context, cmd *cli.Command) error {
	cfg, err := ConfigFromContext(ctx)
	if err != nil {
		return fmt.Errorf("getting config: %w", err)
	}

	loader, err := AwsLoaderFromContext(ctx)
	if err != nil {
		return fmt.Errorf("getting AWS loader: %w", err)
	}

	s3Client, err := loader.S3(ctx)
	if err != nil {
		return fmt.Errorf("creating S3 client: %w", err)
	}

	public := cmd.Bool("public")
	unencrypted := cmd.Bool("unencrypted")
	versioning := cmd.Bool("versioning")

	if !public && !unencrypted && !versioning {
		public = true
		unencrypted = true
		versioning = true
	}

	regions := []string{cfg.AWS.Region}
	if cfg.AWS.AllRegions {
		regions, err = loader.Regions(ctx)
		if err != nil {
			return fmt.Errorf("getting regions: %w", err)
		}
	}

	var results []string

	for _, region := range regions {
		if cfg.AWS.AllRegions {
			fmt.Printf("Checking region: %s\n", region)
		}

		var client *s3.Client
		if cfg.AWS.AllRegions {
			client, err = loader.S3InRegion(ctx, region)
			if err != nil {
				fmt.Printf("Error creating S3 client for region %s: %v\n", region, err)
				continue
			}
		} else {
			client = s3Client
		}

		paginator := s3.NewListBucketsPaginator(client, &s3.ListBucketsInput{})

		for paginator.HasMorePages() {
			output, err := paginator.NextPage(ctx)
			if err != nil {
				fmt.Printf("Error listing buckets in region %s: %v\n", region, err)
				continue
			}

			for _, bucket := range output.Buckets {
				bucketName := *bucket.Name

				if public {
					isPublic, pubErr := checkBucketPublic(ctx, client, bucketName)
					if pubErr != nil {
						if isAccessDenied(pubErr) {
							fmt.Printf("  Warning: Cannot check public access for %s: access denied\n", bucketName)
						}
						continue
					}
					if isPublic {
						results = append(results, fmt.Sprintf("  Bucket: %s - PUBLIC", bucketName))
					}
				}

				if unencrypted {
					isUnencrypted, encErr := checkBucketEncryption(ctx, client, bucketName)
					if encErr != nil {
						if isAccessDenied(encErr) {
							fmt.Printf("  Warning: Cannot check encryption for %s: access denied\n", bucketName)
						}
						continue
					}
					if isUnencrypted {
						results = append(results, fmt.Sprintf("  Bucket: %s - UNENCRYPTED", bucketName))
					}
				}

				if versioning {
					versioningDisabled, verErr := checkBucketVersioning(ctx, client, bucketName)
					if verErr != nil {
						if isAccessDenied(verErr) {
							fmt.Printf("  Warning: Cannot check versioning for %s: access denied\n", bucketName)
						}
						continue
					}
					if versioningDisabled {
						results = append(results, fmt.Sprintf("  Bucket: %s - VERSIONING DISABLED", bucketName))
					}
				}
			}
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

func checkBucketPublic(ctx context.Context, client *s3.Client, bucketName string) (bool, error) {
	acl, err := client.GetBucketAcl(ctx, &s3.GetBucketAclInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return false, err
	}

	for _, grant := range acl.Grants {
		if grant.Grantee.Type == "Group" {
			if grant.Grantee.URI != nil && (*grant.Grantee.URI == "http://acs.amazonaws.com/groups/global/AllUsers" || *grant.Grantee.URI == "http://acs.amazonaws.com/groups/global/AuthenticatedUsers") {
				return true, nil
			}
		}
	}

	return false, nil
}

func checkBucketEncryption(ctx context.Context, client *s3.Client, bucketName string) (bool, error) {
	_, err := client.GetBucketEncryption(ctx, &s3.GetBucketEncryptionInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "ServerSideEncryptionConfigurationNotFound") ||
			strings.Contains(errStr, "The server side encryption configuration was not found") {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func checkBucketVersioning(ctx context.Context, client *s3.Client, bucketName string) (bool, error) {
	output, err := client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return false, err
	}

	return string(output.Status) != "Enabled", nil
}

func isAccessDenied(err error) bool {
	return strings.Contains(err.Error(), "AccessDenied") ||
		strings.Contains(err.Error(), "access denied")
}
