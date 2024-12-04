package ec2

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type VPCs []VPC

func (v VPCs) FilterOut(flowLogs FlowLogs) VPCs {
	var out VPCs
	names := flowLogs.NamesSet()
	for _, vpc := range v {
		if _, ok := names[vpc.Id]; !ok {
			out = append(out, vpc)
		}
	}
	return out
}

type VPC struct {
	Id        string
	Name      string
	Cidr      string
	IsDefault bool
	tags      map[string]string
}

func (v VPC) Tags() map[string]string {
	out := make(map[string]string)
	for k, val := range v.tags {
		out[k] = val
	}
	return out
}

func (v VPC) String() string {
	name := v.Name
	if name == "" {
		name = "-"
	}
	return fmt.Sprintf("%s [%s]", v.Id, name)
}

func toVPCs(in []types.Vpc) VPCs {
	var out VPCs
	for _, v := range in {
		out = append(out, toVPC(v))
	}
	return out
}

func toVPC(in types.Vpc) VPC {
	tags := toTags(in.Tags)
	return VPC{
		Id:        aws.ToString(in.VpcId),
		Name:      tags["Name"],
		Cidr:      aws.ToString(in.CidrBlock),
		IsDefault: aws.ToBool(in.IsDefault),
		tags:      tags,
	}
}
