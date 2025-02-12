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
	OwnerId     string
	GroupName   string
	Description string
	Ingress     []IpPermission
	Egress      []IpPermission
	tags        map[string]string
}

type IpPermission struct {
	FromPort      int
	IpProtocol    string
	IpRanges      []IpRange
	Ipv6Ranges    []IpRange
	PrefixListIds []IdDescription
	ToPort        int
	GroupIds      []IdDescription
}

type IpRange struct {
	Cidr        string
	Description string
}

type IdDescription struct {
	Id          string
	Description string
}

func toIpPermissions(in []types.IpPermission) []IpPermission {
	var out []IpPermission
	for _, v := range in {
		out = append(out, toIpPermission(v))
	}
	return out
}

func toIpPermission(in types.IpPermission) IpPermission {
	var ipRanges []IpRange
	for _, v := range in.IpRanges {
		ipRanges = append(ipRanges, IpRange{
			Cidr:        aws.ToString(v.CidrIp),
			Description: aws.ToString(v.Description),
		})
	}
	var ipv6Ranges []IpRange
	for _, v := range in.Ipv6Ranges {
		ipv6Ranges = append(ipv6Ranges, IpRange{
			Cidr:        aws.ToString(v.CidrIpv6),
			Description: aws.ToString(v.Description),
		})
	}
	var prefixListIds []IdDescription
	for _, v := range in.PrefixListIds {
		prefixListIds = append(prefixListIds, IdDescription{
			Id:          aws.ToString(v.PrefixListId),
			Description: aws.ToString(v.Description),
		})
	}
	var groupIds []IdDescription
	for _, v := range in.UserIdGroupPairs {
		groupIds = append(groupIds, IdDescription{
			Id:          aws.ToString(v.GroupId),
			Description: aws.ToString(v.Description),
		})
	}

	return IpPermission{
		FromPort:      int(aws.ToInt32(in.FromPort)),
		IpProtocol:    aws.ToString(in.IpProtocol),
		IpRanges:      ipRanges,
		Ipv6Ranges:    ipv6Ranges,
		PrefixListIds: prefixListIds,
		ToPort:        int(aws.ToInt32(in.ToPort)),
		GroupIds:      groupIds,
	}
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
		OwnerId:     aws.ToString(in.OwnerId),
		GroupName:   aws.ToString(in.GroupName),
		Description: aws.ToString(in.Description),
		Ingress:     toIpPermissions(in.IpPermissions),
		Egress:      toIpPermissions(in.IpPermissionsEgress),
		tags:        toTags(in.Tags),
	}
}
