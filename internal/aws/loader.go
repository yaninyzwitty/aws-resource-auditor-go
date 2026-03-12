package aws

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	awsconfigpkg "github.com/yaninyzwitty/aws-resource-auditor-go/internal/config"
)

type Loader struct {
	cfg        awsconfigpkg.AWSConfig
	awsCfg     aws.Config
	configOnce sync.Once
	configErr  error

	ec2Once       sync.Once
	ec2Client     *ec2.Client
	s3Once        sync.Once
	s3Client      *s3.Client
	iamOnce       sync.Once
	iamClient     *iam.Client
	lambdaOnce    sync.Once
	lambdaClient  *lambda.Client
	rdsOnce       sync.Once
	rdsClient     *rds.Client
	secretsOnce   sync.Once
	secretsClient *secretsmanager.Client
}

func NewLoader(cfg awsconfigpkg.AWSConfig) (*Loader, error) {
	return &Loader{cfg: cfg}, nil
}

func (l *Loader) loadConfig(ctx context.Context) error {
	l.configOnce.Do(func() {
		var opts []func(*config.LoadOptions) error

		if l.cfg.Profile != "" {
			opts = append(opts, config.WithSharedConfigProfile(l.cfg.Profile))
		}

		cfg, err := config.LoadDefaultConfig(ctx, opts...)
		if err != nil {
			l.configErr = fmt.Errorf("loading AWS config: %w", err)
			return
		}

		if l.cfg.RoleARN != "" {
			provider := stscreds.NewAssumeRoleProvider(
				sts.NewFromConfig(cfg), l.cfg.RoleARN, func(aro *stscreds.AssumeRoleOptions) {
					if l.cfg.ExternalID != "" {
						aro.ExternalID = &l.cfg.ExternalID
					}
				})
			cfg.Credentials = aws.NewCredentialsCache(provider)
		}

		l.awsCfg = cfg
	})

	return l.configErr
}

func (l *Loader) Config(ctx context.Context) (aws.Config, error) {
	if err := l.loadConfig(ctx); err != nil {
		return aws.Config{}, err
	}
	return l.awsCfg, nil
}

func (l *Loader) ConfigForRegion(ctx context.Context, region string) (aws.Config, error) {
	if err := l.loadConfig(ctx); err != nil {
		return aws.Config{}, err
	}

	opts := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}

	if l.cfg.Profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(l.cfg.Profile))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return aws.Config{}, fmt.Errorf("loading AWS config for region %s: %w", region, err)
	}

	if l.cfg.RoleARN != "" {
		provider := stscreds.NewAssumeRoleProvider(
			sts.NewFromConfig(cfg), l.cfg.RoleARN, func(aro *stscreds.AssumeRoleOptions) {
				if l.cfg.ExternalID != "" {
					aro.ExternalID = &l.cfg.ExternalID
				}
			})
		cfg.Credentials = aws.NewCredentialsCache(provider)
	}

	return cfg, nil
}

func (l *Loader) EC2(ctx context.Context) (*ec2.Client, error) {
	if err := l.loadConfig(ctx); err != nil {
		return nil, fmt.Errorf("creating EC2 client: %w", err)
	}

	l.ec2Once.Do(func() {
		l.ec2Client = ec2.NewFromConfig(l.awsCfg)
	})

	return l.ec2Client, nil
}

func (l *Loader) EC2InRegion(ctx context.Context, region string) (*ec2.Client, error) {
	cfg, err := l.ConfigForRegion(ctx, region)
	if err != nil {
		return nil, fmt.Errorf("creating EC2 client for region %s: %w", region, err)
	}
	return ec2.NewFromConfig(cfg), nil
}

func (l *Loader) S3(ctx context.Context) (*s3.Client, error) {
	if err := l.loadConfig(ctx); err != nil {
		return nil, fmt.Errorf("creating S3 client: %w", err)
	}

	l.s3Once.Do(func() {
		l.s3Client = s3.NewFromConfig(l.awsCfg)
	})

	return l.s3Client, nil
}

func (l *Loader) S3InRegion(ctx context.Context, region string) (*s3.Client, error) {
	cfg, err := l.ConfigForRegion(ctx, region)
	if err != nil {
		return nil, fmt.Errorf("creating S3 client for region %s: %w", region, err)
	}
	return s3.NewFromConfig(cfg), nil
}

func (l *Loader) IAM(ctx context.Context) (*iam.Client, error) {
	if err := l.loadConfig(ctx); err != nil {
		return nil, fmt.Errorf("creating IAM client: %w", err)
	}

	l.iamOnce.Do(func() {
		l.iamClient = iam.NewFromConfig(l.awsCfg)
	})

	return l.iamClient, nil
}

func (l *Loader) Lambda(ctx context.Context) (*lambda.Client, error) {
	if err := l.loadConfig(ctx); err != nil {
		return nil, fmt.Errorf("creating Lambda client: %w", err)
	}

	l.lambdaOnce.Do(func() {
		l.lambdaClient = lambda.NewFromConfig(l.awsCfg)
	})

	return l.lambdaClient, nil
}

func (l *Loader) LambdaInRegion(ctx context.Context, region string) (*lambda.Client, error) {
	cfg, err := l.ConfigForRegion(ctx, region)
	if err != nil {
		return nil, fmt.Errorf("creating Lambda client for region %s: %w", region, err)
	}
	return lambda.NewFromConfig(cfg), nil
}

func (l *Loader) RDS(ctx context.Context) (*rds.Client, error) {
	if err := l.loadConfig(ctx); err != nil {
		return nil, fmt.Errorf("creating RDS client: %w", err)
	}

	l.rdsOnce.Do(func() {
		l.rdsClient = rds.NewFromConfig(l.awsCfg)
	})

	return l.rdsClient, nil
}

func (l *Loader) RDSInRegion(ctx context.Context, region string) (*rds.Client, error) {
	cfg, err := l.ConfigForRegion(ctx, region)
	if err != nil {
		return nil, fmt.Errorf("creating RDS client for region %s: %w", region, err)
	}
	return rds.NewFromConfig(cfg), nil
}

func (l *Loader) SecretsManager(ctx context.Context) (*secretsmanager.Client, error) {
	if err := l.loadConfig(ctx); err != nil {
		return nil, fmt.Errorf("creating Secrets Manager client: %w", err)
	}

	l.secretsOnce.Do(func() {
		l.secretsClient = secretsmanager.NewFromConfig(l.awsCfg)
	})

	return l.secretsClient, nil
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

	regions := make([]string, 0, len(resp.Regions))
	for _, r := range resp.Regions {
		if r.RegionName != nil {
			regions = append(regions, *r.RegionName)
		}
	}

	return regions, nil
}
