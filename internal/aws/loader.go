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

func (l *Loader) initConfig(ctx context.Context) error {
	l.configOnce.Do(func() {
		var opts []func(*config.LoadOptions) error

		// Only set region if explicitly provided - let SDK resolve from env/profile otherwise
		if l.cfg.Region != "" {
			opts = append(opts, config.WithRegion(l.cfg.Region))
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
			l.configErr = fmt.Errorf("loading AWS config: %w", err)
			return
		}

		// Validate that region was resolved (unless AllRegions is set, we need a starting region)
		if cfg.Region == "" && !l.cfg.AllRegions {
			l.configErr = fmt.Errorf("region not specified: set AWS_REGION, provide region in config, or use a profile with region")
			return
		}

		l.awsCfg = cfg
	})

	return l.configErr
}

func (l *Loader) Config(ctx context.Context) (aws.Config, error) {
	if err := l.initConfig(ctx); err != nil {
		return aws.Config{}, err
	}
	return l.awsCfg, nil
}

func (l *Loader) EC2(ctx context.Context) (*ec2.Client, error) {
	if err := l.initConfig(ctx); err != nil {
		return nil, fmt.Errorf("creating EC2 client: %w", err)
	}

	l.ec2Once.Do(func() {
		l.ec2Client = ec2.NewFromConfig(l.awsCfg)
	})

	return l.ec2Client, nil
}

func (l *Loader) S3(ctx context.Context) (*s3.Client, error) {
	if err := l.initConfig(ctx); err != nil {
		return nil, fmt.Errorf("creating S3 client: %w", err)
	}

	l.s3Once.Do(func() {
		l.s3Client = s3.NewFromConfig(l.awsCfg)
	})

	return l.s3Client, nil
}

func (l *Loader) IAM(ctx context.Context) (*iam.Client, error) {
	if err := l.initConfig(ctx); err != nil {
		return nil, fmt.Errorf("creating IAM client: %w", err)
	}

	l.iamOnce.Do(func() {
		l.iamClient = iam.NewFromConfig(l.awsCfg)
	})

	return l.iamClient, nil
}

func (l *Loader) Lambda(ctx context.Context) (*lambda.Client, error) {
	if err := l.initConfig(ctx); err != nil {
		return nil, fmt.Errorf("creating Lambda client: %w", err)
	}

	l.lambdaOnce.Do(func() {
		l.lambdaClient = lambda.NewFromConfig(l.awsCfg)
	})

	return l.lambdaClient, nil
}

func (l *Loader) RDS(ctx context.Context) (*rds.Client, error) {
	if err := l.initConfig(ctx); err != nil {
		return nil, fmt.Errorf("creating RDS client: %w", err)
	}

	l.rdsOnce.Do(func() {
		l.rdsClient = rds.NewFromConfig(l.awsCfg)
	})

	return l.rdsClient, nil
}

func (l *Loader) SecretsManager(ctx context.Context) (*secretsmanager.Client, error) {
	if err := l.initConfig(ctx); err != nil {
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

	regions := make([]string, len(resp.Regions))
	for i, r := range resp.Regions {
		regions[i] = *r.RegionName
	}

	return regions, nil
}
