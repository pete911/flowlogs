package ec2

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"slices"
	"time"
)

type FlowLogs []FlowLog

func (f FlowLogs) Ids() []string {
	var out []string
	for _, v := range f {
		out = append(out, v.FlowLogId)
	}
	return out
}

func (f FlowLogs) NamesSet() map[string]struct{} {
	out := make(map[string]struct{})
	for _, v := range f {
		out[v.Name] = struct{}{}
	}
	return out
}

// GetByNames returns sorted names (keys) and map where key is name and value flow logs
func (f FlowLogs) GetByNames() ([]string, map[string]FlowLogs) {
	var names []string
	out := make(map[string]FlowLogs)
	for _, flowLog := range f {
		name := flowLog.Name
		if _, ok := out[name]; !ok {
			out[name] = FlowLogs{}
			names = append(names, name)
		}
		out[name] = append(out[name], flowLog)
	}
	slices.Sort(names)
	return names, out
}

func toFlowLogs(in []types.FlowLog) FlowLogs {
	var out FlowLogs
	for _, v := range in {
		out = append(out, toFlowLog(v))
	}
	return out
}

type FlowLog struct {
	Name         string
	FlowLogId    string
	ResourceId   string
	LogGroupName string
	CreationTime time.Time
	tags         map[string]string
}

func (f FlowLog) Tags() map[string]string {
	out := make(map[string]string)
	for k, v := range f.tags {
		out[k] = v
	}
	return out
}

func toFlowLog(in types.FlowLog) FlowLog {
	tags := toTags(in.Tags)
	return FlowLog{
		Name:         tags["Name"],
		FlowLogId:    aws.ToString(in.FlowLogId),
		ResourceId:   aws.ToString(in.ResourceId),
		LogGroupName: aws.ToString(in.LogGroupName),
		CreationTime: aws.ToTime(in.CreationTime),
		tags:         tags,
	}
}

func (f FlowLog) String() string {
	return fmt.Sprintf("%s [%s - %s]", f.Name, f.FlowLogId, f.ResourceId)
}
