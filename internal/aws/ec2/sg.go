package ec2

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type SecurityGroups []SecurityGroup

func (s SecurityGroups) FilterOut(flowLogs FlowLogs) SecurityGroups {
	var out SecurityGroups
	names := flowLogs.NamesSet()
	for _, v := range s {
		if _, ok := names[v.Id]; !ok {
			out = append(out, v)
		}
	}
	return out
}

type SecurityGroup struct {
	VpcId       string
	Id          string
	GroupName   string
	Description string
	tags        map[string]string
}

func (s SecurityGroup) Tags() map[string]string {
	out := make(map[string]string)
	for k, v := range s.tags {
		out[k] = v
	}
	return out
}

func (s SecurityGroup) String() string {
	name := s.GroupName
	if name == "" {
		name = "-"
	}
	return fmt.Sprintf("%s [%s]", s.Id, name)
}

func toSecurityGroups(in []types.SecurityGroup) SecurityGroups {
	var out SecurityGroups
	for _, v := range in {
		out = append(out, toSecurityGroup(v))
	}
	return out
}

func toSecurityGroup(in types.SecurityGroup) SecurityGroup {
	return SecurityGroup{
		VpcId:       aws.ToString(in.VpcId),
		Id:          aws.ToString(in.GroupId),
		GroupName:   aws.ToString(in.GroupName),
		Description: aws.ToString(in.Description),
		tags:        toTags(in.Tags),
	}
}
