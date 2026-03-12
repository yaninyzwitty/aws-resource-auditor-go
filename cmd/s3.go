package cmd

import (
	"context"
	"fmt"

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

	var results []string

	paginator := s3.NewListBucketsPaginator(s3Client, &s3.ListBucketsInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("listing buckets: %w", err)
		}

		for _, bucket := range output.Buckets {
			bucketName := *bucket.Name

			if public {
				isPublic, err := checkBucketPublic(ctx, s3Client, bucketName)
				if err == nil && isPublic {
					results = append(results, fmt.Sprintf("  Bucket: %s - PUBLIC", bucketName))
				}
			}

			if unencrypted {
				hasEncryption, err := checkBucketEncryption(ctx, s3Client, bucketName)
				if err == nil && !hasEncryption {
					results = append(results, fmt.Sprintf("  Bucket: %s - UNENCRYPTED", bucketName))
				}
			}

			if versioning {
				versioningEnabled, err := checkBucketVersioning(ctx, s3Client, bucketName)
				if err == nil && !versioningEnabled {
					results = append(results, fmt.Sprintf("  Bucket: %s - VERSIONING DISABLED", bucketName))
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
		return false, err
	}
	return true, nil
}

func checkBucketVersioning(ctx context.Context, client *s3.Client, bucketName string) (bool, error) {
	output, err := client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return false, err
	}

	return string(output.Status) == "Enabled", nil
}
