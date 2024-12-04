package ec2

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type Subnets []Subnet

func (s Subnets) FilterOut(flowLogs FlowLogs) Subnets {
	var out Subnets
	names := flowLogs.NamesSet()
	for _, v := range s {
		if _, ok := names[v.Id]; !ok {
			out = append(out, v)
		}
	}
	return out
}

type Subnet struct {
	VpcId                   string
	Id                      string
	SubnetArn               string
	Name                    string
	AvailabilityZone        string
	AvailabilityZoneId      string
	AvailableIpAddressCount int
	CidrBlock               string
	DefaultForAz            bool
	tags                    map[string]string
}

func (s Subnet) Tags() map[string]string {
	out := make(map[string]string)
	for k, v := range s.tags {
		out[k] = v
	}
	return out
}

func (s Subnet) String() string {
	name := s.Name
	if name == "" {
		name = "-"
	}
	return fmt.Sprintf("%s [%s]", s.Id, name)
}

func toSubnets(in []types.Subnet) Subnets {
	var out Subnets
	for _, v := range in {
		out = append(out, toSubnet(v))
	}
	return out
}

func toSubnet(in types.Subnet) Subnet {
	tags := toTags(in.Tags)
	return Subnet{
		VpcId:                   aws.ToString(in.VpcId),
		Id:                      aws.ToString(in.SubnetId),
		SubnetArn:               aws.ToString(in.SubnetArn),
		Name:                    tags["Name"],
		AvailabilityZone:        aws.ToString(in.AvailabilityZone),
		AvailabilityZoneId:      aws.ToString(in.AvailabilityZoneId),
		AvailableIpAddressCount: int(aws.ToInt32(in.AvailableIpAddressCount)),
		CidrBlock:               aws.ToString(in.CidrBlock),
		DefaultForAz:            aws.ToBool(in.DefaultForAz),
		tags:                    tags,
	}
}
