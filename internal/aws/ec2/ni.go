package ec2

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"strings"
	"time"
)

type NetworkInterfaces []NetworkInterface

type NetworkInterface struct {
	VpcId              string
	SubnetId           string
	OwnerId            string
	PublicIP           string
	PublicDnsName      string
	PrivateIpAddress   string
	PrivateIpAddresses []string
	PrivateDnsName     string
	AttachTime         time.Time
	AvailabilityZone   string
	Description        string
	InterfaceType      string
	NetworkInterfaceId string
	RequesterId        string
	RequesterManaged   bool
	InstanceId         string
	Type               string
	Name               string
	Status             string
}

func ToNetworkInterfaces(in []types.NetworkInterface) NetworkInterfaces {
	var out NetworkInterfaces
	for _, v := range in {
		out = append(out, ToNetworkInterface(v))
	}
	return out
}

func ToNetworkInterface(in types.NetworkInterface) NetworkInterface {
	var publicIp, publicDnsName string
	if in.Association != nil {
		publicIp = aws.ToString(in.Association.PublicIp)
		publicDnsName = aws.ToString(in.Association.PublicDnsName)
	}

	var attachTime time.Time
	if in.Attachment != nil {
		attachTime = aws.ToTime(in.Attachment.AttachTime)
	}

	var instanceId string
	if in.Attachment != nil {
		instanceId = aws.ToString(in.Attachment.InstanceId)
	}

	var privateIpAddresses []string
	for _, address := range in.PrivateIpAddresses {
		privateIpAddresses = append(privateIpAddresses, aws.ToString(address.PrivateIpAddress))
	}

	vpcId := aws.ToString(in.VpcId)
	return NetworkInterface{
		VpcId:              vpcId,
		SubnetId:           aws.ToString(in.SubnetId),
		OwnerId:            aws.ToString(in.OwnerId),
		PublicIP:           publicIp,
		PublicDnsName:      publicDnsName,
		PrivateIpAddress:   aws.ToString(in.PrivateIpAddress),
		PrivateIpAddresses: privateIpAddresses,
		PrivateDnsName:     aws.ToString(in.PrivateDnsName),
		AttachTime:         attachTime,
		AvailabilityZone:   aws.ToString(in.AvailabilityZone),
		Description:        aws.ToString(in.Description),
		InterfaceType:      string(in.InterfaceType),
		NetworkInterfaceId: aws.ToString(in.NetworkInterfaceId),
		RequesterId:        aws.ToString(in.RequesterId),
		RequesterManaged:   aws.ToBool(in.RequesterManaged),
		InstanceId:         instanceId,
		Type:               getNiType(in),
		Name:               getNiName(in),
		Status:             string(in.Status),
	}
}

func (v NetworkInterfaces) GetById(id string) NetworkInterface {
	for _, ni := range v {
		if ni.NetworkInterfaceId == id {
			return ni
		}
	}
	return NetworkInterface{}
}

func (v NetworkInterfaces) GetByIp(ip string) NetworkInterfaces {
	matcher := func(in string) bool {
		return in == ip
	}

	var out NetworkInterfaces
	for _, ni := range v {
		if ni.matchesIp(matcher) {
			out = append(out, ni)
			continue
		}
	}
	return out
}

func (v NetworkInterface) matchesIp(matcher func(in string) bool) bool {
	// private ip address is already in private ip addresses slice, but just in case check all
	for _, ip := range append(v.PrivateIpAddresses, v.PrivateIpAddress, v.PublicIP) {
		if matcher(ip) {
			return true
		}
	}
	return false
}

func getNiName(in types.NetworkInterface) string {
	if in.Attachment != nil {
		if id := aws.ToString(in.Attachment.InstanceId); id != "" {
			return id
		}
	}

	requesterId := aws.ToString(in.RequesterId)
	description := aws.ToString(in.Description)

	// application or network load balancers
	if strings.HasPrefix(description, "ELB app/") || strings.HasPrefix(description, "ELB net/") {
		parts := strings.Split(description, "/")
		if len(parts) >= 2 {
			return parts[1]
		}
		return ""
	}

	// classic load balancer has requester id amazon-elb and description prefix 'ELB '
	if requesterId == "amazon-elb" && strings.HasPrefix(description, "ELB ") {
		return strings.TrimPrefix(description, "ELB ")
	}

	if strings.HasPrefix(description, "ElastiCache ") {
		return strings.TrimPrefix(description, "ElastiCache ")
	}
	if strings.HasPrefix(description, "ElastiCache+") {
		return strings.TrimPrefix(description, "ElastiCache+")
	}

	if strings.HasPrefix(description, "AWS Lambda VPC ENI-") {
		nameWithSuffix := strings.TrimPrefix(description, "AWS Lambda VPC ENI-")
		// remove xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx suffix
		parts := strings.Split(nameWithSuffix, "-")
		if len(parts) > 5 {
			return strings.Join(parts[:len(parts)-5], "-")
		}
		return nameWithSuffix
	}

	if strings.HasPrefix(description, "datasync ") {
		return strings.TrimPrefix(description, "datasync ")
	}
	if strings.HasPrefix(description, "Interface for NAT Gateway nat-") {
		return strings.TrimPrefix(description, "Interface for NAT Gateway ")
	}
	return "-"
}

func getNiType(in types.NetworkInterface) string {
	if in.Attachment != nil {
		if aws.ToString(in.Attachment.InstanceId) != "" {
			return "instance"
		}
	}

	requesterId := aws.ToString(in.RequesterId)
	description := aws.ToString(in.Description)
	// network load balancer has requester id account id and description prefix 'ELB net/', but the interface type
	// is set to network_load_balancer as well
	if in.InterfaceType == types.NetworkInterfaceTypeNetworkLoadBalancer {
		return "nlb"
	}
	// application load balancer has requester id amazon-elb and description prefix 'ELB app/'
	if strings.HasPrefix(description, "ELB app/") {
		return "alb"
	}
	// classic load balancer has requester id amazon-elb and description prefix 'ELB '
	if requesterId == "amazon-elb" && strings.HasPrefix(description, "ELB ") {
		return "elb"
	}

	if strings.HasPrefix(description, "ElastiCache ") {
		return "elastic_cache"
	}
	if strings.HasPrefix(description, "AWS Lambda VPC ENI-") {
		return "lambda"
	}
	if strings.HasPrefix(description, "datasync ") {
		return "datasync"
	}
	if strings.HasPrefix(description, "[Do not delete] Network Interface created to access resources in your VPC for SageMaker Notebook Instance ") {
		return "sage_maker"
	}
	if strings.HasPrefix(description, "[DO NOT DELETE] ENI managed by SageMaker for Studio Domain") {
		return "sage_maker"
	}
	if strings.HasPrefix(description, "Attached to Glue using role: arn:aws:iam::") {
		return "glue"
	}
	if strings.HasPrefix(description, "arn:aws:ecs:") {
		return "ecs"
	}
	if strings.HasPrefix(description, "AWS created network interface for directory ") {
		return "ad"
	}
	if strings.HasPrefix(description, "Created By Amazon Workspaces for AWS Account ID ") {
		return "workspace"
	}
	if description == "RDSNetworkInterface" {
		return "rds"
	}
	if description == "RedshiftNetworkInterface" {
		return "redshift"
	}
	if in.InterfaceType == types.NetworkInterfaceTypeNatGateway {
		return "nat"
	}
	if strings.HasPrefix(description, "Interface for NAT Gateway nat-") {
		return "nat"
	}
	return string(in.InterfaceType)
}
