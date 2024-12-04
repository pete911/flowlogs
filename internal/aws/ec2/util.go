package ec2

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func toTags(in []types.Tag) map[string]string {
	out := make(map[string]string)
	for _, tag := range in {
		out[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
	}
	return out
}
