package logs

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"time"
)

type LogGroup struct {
	LogGroupArn     string
	CreationTime    time.Time
	KmsKeyId        string
	LogGroupClass   string
	LogGroupName    string
	RetentionInDays int
	StoredBytes     int
}

func toLogGroup(in types.LogGroup) LogGroup {
	return LogGroup{
		LogGroupArn:     aws.ToString(in.LogGroupArn),
		CreationTime:    time.UnixMilli(aws.ToInt64(in.CreationTime)),
		KmsKeyId:        aws.ToString(in.KmsKeyId),
		LogGroupClass:   string(in.LogGroupClass),
		LogGroupName:    aws.ToString(in.LogGroupName),
		RetentionInDays: int(aws.ToInt32(in.RetentionInDays)),
		StoredBytes:     int(aws.ToInt64(in.StoredBytes)),
	}
}
