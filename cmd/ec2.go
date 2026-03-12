package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
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
		Action: ec2Action,
	}
}

func ec2Action(ctx context.Context, cmd *cli.Command) error {
	cfg, err := ConfigFromContext(ctx)
	if err != nil {
		return fmt.Errorf("getting config: %w", err)
	}

	loader, err := AwsLoaderFromContext(ctx)
	if err != nil {
		return fmt.Errorf("getting AWS loader: %w", err)
	}

	ec2Client, err := loader.EC2(ctx)
	if err != nil {
		return fmt.Errorf("creating EC2 client: %w", err)
	}

	regions := []string{cfg.AWS.Region}
	if cfg.AWS.AllRegions {
		regions, err = loader.Regions(ctx)
		if err != nil {
			return fmt.Errorf("getting regions: %w", err)
		}
	}

	unused := cmd.Bool("unused")
	oldAmis := cmd.Bool("old-amis")
	unencrypted := cmd.Bool("unencrypted")

	if !unused && !oldAmis && !unencrypted {
		unused = true
		oldAmis = true
		unencrypted = true
	}

	var results []string

	for _, region := range regions {
		if cfg.AWS.AllRegions {
			fmt.Printf("Checking region: %s\n", region)
		}

		var client *ec2.Client
		if cfg.AWS.AllRegions {
			client, err = loader.EC2InRegion(ctx, region)
			if err != nil {
				fmt.Printf("Error creating EC2 client for region %s: %v\n", region, err)
				continue
			}
		} else {
			client = ec2Client
		}

		if unused {
			instances, err := findUnusedInstances(ctx, client, cfg.Thresholds.OlderThan)
			if err != nil {
				fmt.Printf("Error finding unused instances: %v\n", err)
			}
			results = append(results, instances...)
		}

		if oldAmis {
			amis, err := findOldAMIs(ctx, client, cfg.Services.EC2.OldAMIsOlderThan)
			if err != nil {
				fmt.Printf("Error finding old AMIs: %v\n", err)
			}
			results = append(results, amis...)
		}

		if unencrypted {
			volumes, err := findUnencryptedVolumes(ctx, client)
			if err != nil {
				fmt.Printf("Error finding unencrypted volumes: %v\n", err)
			}
			results = append(results, volumes...)
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

func findUnusedInstances(ctx context.Context, client *ec2.Client, olderThan time.Duration) ([]string, error) {
	var results []string

	paginator := ec2.NewDescribeInstancesPaginator(client, &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{Name: aws.String("instance-state-name"), Values: []string{"running", "stopped"}},
		},
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return results, fmt.Errorf("describing instances: %w", err)
		}

		for _, reservation := range output.Reservations {
			for _, instance := range reservation.Instances {
				if instance.LaunchTime == nil {
					continue
				}
				age := time.Since(*instance.LaunchTime)

				if olderThan > 0 && age < olderThan {
					continue
				}

				var name string
				for _, tag := range instance.Tags {
					if tag.Key != nil && *tag.Key == "Name" && tag.Value != nil {
						name = *tag.Value
						break
					}
				}

				result := fmt.Sprintf("  Instance: %s (%s) - State: %s - Launched: %s ago",
					*instance.InstanceId,
					name,
					string(instance.State.Name),
					formatDuration(age),
				)
				results = append(results, result)
			}
		}
	}

	return results, nil
}

func findOldAMIs(ctx context.Context, client *ec2.Client, olderThan time.Duration) ([]string, error) {
	var results []string

	if olderThan == 0 {
		olderThan = 90 * 24 * time.Hour
	}

	paginator := ec2.NewDescribeImagesPaginator(client, &ec2.DescribeImagesInput{
		Owners: []string{"self"},
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return results, fmt.Errorf("describing images: %w", err)
		}

		for _, image := range output.Images {
			if image.CreationDate == nil {
				continue
			}
			creationTime, err := time.Parse(time.RFC3339, *image.CreationDate)
			if err != nil {
				continue
			}
			age := time.Since(creationTime)

			if age < olderThan {
				continue
			}

			imageName := ""
			if image.Name != nil {
				imageName = *image.Name
			}

			result := fmt.Sprintf("  AMI: %s (%s) - Created: %s ago",
				*image.ImageId,
				imageName,
				formatDuration(age),
			)
			results = append(results, result)
		}
	}

	return results, nil
}

func findUnencryptedVolumes(ctx context.Context, client *ec2.Client) ([]string, error) {
	var results []string

	paginator := ec2.NewDescribeVolumesPaginator(client, &ec2.DescribeVolumesInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return results, fmt.Errorf("describing volumes: %w", err)
		}

		for _, volume := range output.Volumes {
			if volume.Encrypted != nil && *volume.Encrypted {
				continue
			}

			result := fmt.Sprintf("  Volume: %s - Size: %d GB - State: %s",
				*volume.VolumeId,
				*volume.Size,
				string(volume.State),
			)
			results = append(results, result)
		}
	}

	return results, nil
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	if days > 0 {
		return fmt.Sprintf("%dd", days)
	}
	hours := int(d.Hours())
	if hours > 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dm", int(d.Minutes()))
}
