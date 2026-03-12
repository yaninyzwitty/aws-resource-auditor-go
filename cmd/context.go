package cmd

import (
	"context"
	"errors"

	"github.com/yaninyzwitty/aws-resource-auditor-go/internal/aws"
	"github.com/yaninyzwitty/aws-resource-auditor-go/internal/config"
)

type contextKey string

const (
	ConfigKey    contextKey = "config"
	AwsLoaderKey contextKey = "awsLoader"
)

type Config = config.Config
type AwsLoader = aws.Loader

var (
	ErrConfigNotFound    = errors.New("config not found in context")
	ErrAwsLoaderNotFound = errors.New("AWS loader not found in context")
)

func ConfigFromContext(ctx context.Context) (*Config, error) {
	cfg, ok := ctx.Value(ConfigKey).(*Config)
	if !ok || cfg == nil {
		return nil, ErrConfigNotFound
	}
	return cfg, nil
}

func AwsLoaderFromContext(ctx context.Context) (*AwsLoader, error) {
	loader, ok := ctx.Value(AwsLoaderKey).(*AwsLoader)
	if !ok || loader == nil {
		return nil, ErrAwsLoaderNotFound
	}
	return loader, nil
}
