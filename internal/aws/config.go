package aws

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"time"
)

type Config struct {
	Account string
	Region  string
	Config  aws.Config
}

func NewConfig(awsRegion string) (Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return Config{}, fmt.Errorf("load aws config: %w", err)
	}

	if awsRegion == "" && cfg.Region == "" {
		return Config{}, errors.New("missing aws region")
	}

	if awsRegion != "" {
		cfg.Region = awsRegion
	}

	account, err := getCurrentAWSAccount(cfg)
	if err != nil {
		return Config{}, fmt.Errorf("get current aws account: %w", err)
	}

	return Config{
		Account: account,
		Region:  cfg.Region,
		Config:  cfg,
	}, nil
}

func getCurrentAWSAccount(cfg aws.Config) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	svc := sts.NewFromConfig(cfg)
	resp, err := svc.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}
	return aws.ToString(resp.Account), nil
}
