package iam

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"log/slog"
	"time"
)

const (
	description = "flowlogs cli role"
	policyName  = "flow-logs"

	trustPolicy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "vpc-flow-logs.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}`
	policy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogStream",
        "logs:PutLogEvents",
        "logs:DescribeLogGroups",
        "logs:DescribeLogStreams"
      ],
      "Resource": "*"
    }
  ]
}`
)

type Client struct {
	logger *slog.Logger
	svc    *iam.Client
}

func NewClient(logger *slog.Logger, cfg aws.Config) Client {
	return Client{
		logger: logger,
		svc:    iam.NewFromConfig(cfg),
	}
}

// CreateFlowLogsRole create role and return arn
func (c Client) CreateFlowLogsRole(roleName string, tags map[string]string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	in := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(trustPolicy),
		RoleName:                 aws.String(roleName),
		Description:              aws.String(description),
		Tags:                     fromTags(tags),
	}

	out, err := c.svc.CreateRole(ctx, in)
	if err != nil {
		return "", err
	}
	if err := c.putRolePolicy(roleName, policyName, policy); err != nil {
		return "", fmt.Errorf("put %s role policy: %w", roleName, err)
	}
	return aws.ToString(out.Role.Arn), nil
}

func (c Client) putRolePolicy(roleName, policyName, policyDocument string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	in := &iam.PutRolePolicyInput{
		PolicyDocument: aws.String(policyDocument),
		PolicyName:     aws.String(policyName),
		RoleName:       aws.String(roleName),
	}

	if _, err := c.svc.PutRolePolicy(ctx, in); err != nil {
		return err
	}
	return nil
}

func (c Client) DeleteRole(roleName string, tags map[string]string) error {
	// make sure we are deleting the right role, check tags before deletion
	if err := c.roleMatches(roleName, tags); err != nil {
		return err
	}
	c.logger.Debug(fmt.Sprintf("role %s matches supplied tags", roleName))

	// delete policy first, otherwise we get 'DeleteConflict: Cannot delete entity, must delete policies first' error
	if err := c.deleteRolePolicy(roleName, policyName); err != nil {
		return err
	}
	c.logger.Debug(fmt.Sprintf("role %s policy %s deleted", roleName, policyName))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	in := &iam.DeleteRoleInput{RoleName: aws.String(roleName)}

	if _, err := c.svc.DeleteRole(ctx, in); err != nil {
		return err
	}
	return nil
}

func (c Client) deleteRolePolicy(roleName, policyName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	in := &iam.DeleteRolePolicyInput{
		PolicyName: aws.String(policyName),
		RoleName:   aws.String(roleName),
	}

	if _, err := c.svc.DeleteRolePolicy(ctx, in); err != nil {
		return err
	}
	return nil
}

func (c Client) roleMatches(roleName string, tags map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	in := &iam.GetRoleInput{RoleName: aws.String(roleName)}

	out, err := c.svc.GetRole(ctx, in)
	if err != nil {
		return err
	}

	// we are checking if the role matches supplied tags, not the other way around. Meaning, that if the role has
	// additional tag(s), role matches
	// we can add validation for policy and trust policy as well if needed
	roleTags := toTags(out.Role.Tags)
	for k, v := range tags {
		roleTagValue, ok := roleTags[k]
		if !ok {
			return fmt.Errorf("role %s does not have %s key", roleName, k)
		}
		if roleTagValue != v {
			return fmt.Errorf("role %s key %s value %s does not match %s", roleName, k, roleTagValue, v)
		}
	}
	return nil
}

func fromTags(in map[string]string) []types.Tag {
	var out []types.Tag
	for k, v := range in {
		out = append(out, types.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	return out
}
func toTags(in []types.Tag) map[string]string {
	out := make(map[string]string)
	for _, tag := range in {
		out[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
	}
	return out
}
