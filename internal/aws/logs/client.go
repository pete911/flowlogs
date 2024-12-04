package logs

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"log/slog"
	"time"
)

const retentionDays = 30

type Client struct {
	logger *slog.Logger
	svc    *cloudwatchlogs.Client
}

func NewClient(logger *slog.Logger, cfg aws.Config) Client {
	return Client{
		logger: logger,
		svc:    cloudwatchlogs.NewFromConfig(cfg),
	}
}

func (c Client) CreateLogGroup(logGroupName string, tags map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	in := &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName:  aws.String(logGroupName),
		LogGroupClass: types.LogGroupClassStandard,
		Tags:          tags,
	}

	if _, err := c.svc.CreateLogGroup(ctx, in); err != nil {
		return err
	}

	if err := c.putRetentionPolicy(logGroupName, retentionDays); err != nil {
		return fmt.Errorf("put retention policy on %s log group for %d days: %w", logGroupName, retentionDays, err)
	}
	return nil
}

func (c Client) putRetentionPolicy(logGroupName string, retentionDays int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	in := &cloudwatchlogs.PutRetentionPolicyInput{
		LogGroupName:    aws.String(logGroupName),
		RetentionInDays: aws.Int32(int32(retentionDays)),
	}

	_, err := c.svc.PutRetentionPolicy(ctx, in)
	return err
}

func (c Client) DeleteLogGroup(name string, tags map[string]string) error {
	// make sure we are deleting correct log group, validate tags as well
	if err := c.logGroupMatches(name, tags); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	in := &cloudwatchlogs.DeleteLogGroupInput{LogGroupName: aws.String(name)}

	_, err := c.svc.DeleteLogGroup(ctx, in)
	return err
}

func (c Client) logGroupMatches(logGroupName string, tags map[string]string) error {
	logGroup, err := c.describeLogGroup(logGroupName)
	if err != nil {
		return err
	}

	// list and compare tags
	logGroupTags, err := c.listTags(logGroup.LogGroupArn)
	if err != nil {
		return err
	}

	// we are checking if the log group matches supplied tags, not the other way around. Meaning, that if the log group
	// has additional tag(s), log group matches
	for k, v := range tags {
		logGroupTagValue, ok := logGroupTags[k]
		if !ok {
			return fmt.Errorf("log group %s does not have %s key", logGroupName, k)
		}
		if logGroupTagValue != v {
			return fmt.Errorf("log group %s key %s value %s does not match %s", logGroupName, k, logGroupTagValue, v)
		}
	}
	return nil
}

func (c Client) describeLogGroup(logGroupName string) (LogGroup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	in := &cloudwatchlogs.DescribeLogGroupsInput{LogGroupNamePrefix: aws.String(logGroupName)}

	var logGroups []types.LogGroup
	for {
		out, err := c.svc.DescribeLogGroups(ctx, in)
		if err != nil {
			return LogGroup{}, err
		}
		logGroups = append(logGroups, out.LogGroups...)
		if aws.ToString(out.NextToken) == "" {
			break
		}
		in.NextToken = out.NextToken
	}

	for _, lg := range logGroups {
		if aws.ToString(lg.LogGroupName) == logGroupName {
			return toLogGroup(lg), nil
		}
	}
	return LogGroup{}, fmt.Errorf("log group %s not found", logGroupName)
}

func (c Client) listTags(logGroupArn string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	in := &cloudwatchlogs.ListTagsForResourceInput{ResourceArn: aws.String(logGroupArn)}

	out, err := c.svc.ListTagsForResource(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("log group arn %s: %w", logGroupArn, err)
	}
	return out.Tags, nil
}

func (c Client) ListLogGroups(logGroupNamePrefix string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	in := &cloudwatchlogs.DescribeLogGroupsInput{LogGroupNamePrefix: aws.String(logGroupNamePrefix)}

	var logGroups []string
	for {
		out, err := c.svc.DescribeLogGroups(ctx, in)
		if err != nil {
			return nil, err
		}
		for _, logGroup := range out.LogGroups {
			logGroups = append(logGroups, aws.ToString(logGroup.LogGroupName))
		}
		if aws.ToString(out.NextToken) == "" {
			break
		}
		in.NextToken = out.NextToken
	}
	return logGroups, nil
}

func (c Client) Query(logGroupNames []string, queryString string, sinceMinutes, limit int) ([]map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	in := &cloudwatchlogs.StartQueryInput{
		EndTime:       aws.Int64(time.Now().Unix()),
		StartTime:     aws.Int64(time.Now().Add(time.Duration(-sinceMinutes) * time.Minute).Unix()),
		QueryString:   aws.String(queryString),
		Limit:         aws.Int32(int32(limit)),
		LogGroupNames: logGroupNames,
	}

	out, err := c.svc.StartQuery(ctx, in)
	if err != nil {
		return nil, err
	}
	return c.getQueryResults(aws.ToString(out.QueryId))
}

func (c Client) getQueryResults(queryId string) ([]map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	in := cloudwatchlogs.GetQueryResultsInput{QueryId: aws.String(queryId)}

	// wait before making first call
	retrySecond := 2
	time.Sleep(time.Duration(retrySecond) * time.Second)

	var retry int
	for {
		if retry > 5 {
			return nil, fmt.Errorf("retried %d times, query %s failed", retry, queryId)
		}
		out, err := c.svc.GetQueryResults(ctx, &in)
		if err != nil {
			return nil, err
		}
		// first we check if query is still running (then we retry), or failed (fail fast)
		// Cancelled , Complete , Failed , Running , Scheduled , Timeout , and Unknown .
		if out.Status == types.QueryStatusRunning {
			c.logger.Info(fmt.Sprintf("query status %s, retrying in %d second", out.Status, retrySecond))
			time.Sleep(time.Duration(retrySecond) * time.Second)
			retry++
			continue
		}
		if out.Status != types.QueryStatusComplete {
			return nil, fmt.Errorf("query %s status %s", queryId, out.Status)
		}
		return toQueryResults(out.Results), nil
	}
}

func toQueryResults(in [][]types.ResultField) []map[string]string {
	var out []map[string]string
	for _, line := range in {
		logLine := make(map[string]string)
		for _, field := range line {
			logLine[aws.ToString(field.Field)] = aws.ToString(field.Value)
		}
		out = append(out, logLine)
	}
	return out
}
