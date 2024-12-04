package ec2

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type NatGateways []NatGateway

func (n NatGateways) FilterOut(flowLogs FlowLogs) NatGateways {
	var out NatGateways
	names := flowLogs.NamesSet()
	for _, v := range n {
		if _, ok := names[v.Id]; !ok {
			out = append(out, v)
		}
	}
	return out
}

type NatGateway struct {
	Id                 string
	NetworkInterfaceId string
	SubnetId           string
	Tags               map[string]string
}

func (n NatGateway) String() string {
	name := n.Tags["Name"]
	if name == "" {
		name = "-"
	}
	return fmt.Sprintf("%s [%s]", n.Id, name)
}

func toNatGateways(in []types.NatGateway) NatGateways {
	var out NatGateways
	for _, v := range in {
		out = append(out, toNatGateway(v))
	}
	return out
}

func toNatGateway(in types.NatGateway) NatGateway {
	var networkInterfaceId string
	for _, address := range in.NatGatewayAddresses {
		if aws.ToBool(address.IsPrimary) {
			networkInterfaceId = aws.ToString(address.NetworkInterfaceId)
			break
		}
	}

	return NatGateway{
		Id:                 aws.ToString(in.NatGatewayId),
		NetworkInterfaceId: networkInterfaceId,
		SubnetId:           aws.ToString(in.SubnetId),
		Tags:               toTags(in.Tags),
	}
}
