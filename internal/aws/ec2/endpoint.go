package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type VPCEndpoints []VPCEndpoint

func (v VPCEndpoints) FilterOut(flowLogs FlowLogs) VPCEndpoints {
	var out VPCEndpoints
	names := flowLogs.NamesSet()
	for _, endpoint := range v {
		if _, ok := names[endpoint.VpcEndpointId]; !ok {
			out = append(out, endpoint)
		}
	}
	return out
}

type VPCEndpoint struct {
	VpcEndpointId       string
	VpcEndpointType     string
	DnsNames            []string
	NetworkInterfaceIds []string
	VpcId               string
	SubnetIds           []string
	OwnerId             string
	ServiceName         string
	ServiceNetworkArn   string
	State               string
	tags                map[string]string
}

func (v VPCEndpoint) Tags() map[string]string {
	out := make(map[string]string)
	for k, val := range v.tags {
		out[k] = val
	}
	return out
}

func (v VPCEndpoint) String() string {
	name := v.ServiceName
	if v, ok := v.tags["Name"]; ok {
		name = v + " - " + name
	}
	return fmt.Sprintf("%s [%s] [vpc %s]", v.VpcEndpointId, name, v.VpcId)
}

func toVPCEndpoints(in []types.VpcEndpoint) VPCEndpoints {
	var out VPCEndpoints
	for _, v := range in {
		out = append(out, toVPCEndpoint(v))
	}
	return out
}

func toVPCEndpoint(in types.VpcEndpoint) VPCEndpoint {
	var dnsNames []string
	for _, v := range in.DnsEntries {
		dnsNames = append(dnsNames, aws.ToString(v.DnsName))
	}

	return VPCEndpoint{
		VpcEndpointId:       aws.ToString(in.VpcEndpointId),
		VpcEndpointType:     string(in.VpcEndpointType),
		DnsNames:            dnsNames,
		NetworkInterfaceIds: in.NetworkInterfaceIds,
		VpcId:               aws.ToString(in.VpcId),
		SubnetIds:           in.SubnetIds,
		OwnerId:             aws.ToString(in.OwnerId),
		ServiceName:         aws.ToString(in.ServiceName),
		ServiceNetworkArn:   aws.ToString(in.ServiceNetworkArn),
		State:               string(in.State),
		tags:                toTags(in.Tags),
	}
}
