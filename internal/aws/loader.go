package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	awsconfigpkg "github.com/yaninyzwitty/aws-resource-auditor-go/internal/config"
)

type Loader struct {
	cfg awsconfigpkg.AWSConfig
}

func NewLoader(cfg awsconfigpkg.AWSConfig) (*Loader, error) {
	if cfg.Region == "" && !cfg.AllRegions {
		return nil, fmt.Errorf("region must be specified")
	}

	return &Loader{cfg: cfg}, nil
}

func (l *Loader) Config(ctx context.Context) (aws.Config, error) {
	opts := []func(*config.LoadOptions) error{
		config.WithRegion(l.cfg.Region),
	}

	if l.cfg.Profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(l.cfg.Profile))
	}

	if l.cfg.RoleARN != "" {
		opts = append(opts, config.WithAssumeRoleCredentialOptions(func(o *stscreds.AssumeRoleOptions) {
			o.RoleARN = l.cfg.RoleARN
			if l.cfg.ExternalID != "" {
				o.ExternalID = &l.cfg.ExternalID
			}
		}))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return cfg, fmt.Errorf("loading AWS config: %w", err)
	}

	return cfg, nil
}

func (l *Loader) EC2(ctx context.Context) (*ec2.Client, error) {
	cfg, err := l.Config(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating EC2 client: %w", err)
	}
	return ec2.NewFromConfig(cfg), nil
}

func (l *Loader) S3(ctx context.Context) (*s3.Client, error) {
	cfg, err := l.Config(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating S3 client: %w", err)
	}
	return s3.NewFromConfig(cfg), nil
}

func (l *Loader) IAM(ctx context.Context) (*iam.Client, error) {
	cfg, err := l.Config(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating IAM client: %w", err)
	}
	return iam.NewFromConfig(cfg), nil
}

func (l *Loader) Lambda(ctx context.Context) (*lambda.Client, error) {
	cfg, err := l.Config(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating Lambda client: %w", err)
	}
	return lambda.NewFromConfig(cfg), nil
}

func (l *Loader) RDS(ctx context.Context) (*rds.Client, error) {
	cfg, err := l.Config(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating RDS client: %w", err)
	}
	return rds.NewFromConfig(cfg), nil
}

func (l *Loader) SecretsManager(ctx context.Context) (*secretsmanager.Client, error) {
	cfg, err := l.Config(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating Secrets Manager client: %w", err)
	}
	return secretsmanager.NewFromConfig(cfg), nil
}

func (l *Loader) Regions(ctx context.Context) ([]string, error) {
	ec2Client, err := l.EC2(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating EC2 client for regions: %w", err)
	}

	resp, err := ec2Client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, fmt.Errorf("describing regions: %w", err)
	}

	regions := make([]string, len(resp.Regions))
	for i, r := range resp.Regions {
		regions[i] = *r.RegionName
	}

	return regions, nil
}
