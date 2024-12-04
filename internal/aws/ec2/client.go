package ec2

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/pete911/flowlogs/internal/aws/fields"
	"log/slog"
	"strings"
	"time"
)

type Client struct {
	logger *slog.Logger
	svc    *ec2.Client
}

func NewClient(logger *slog.Logger, cfg aws.Config) Client {
	return Client{
		logger: logger,
		svc:    ec2.NewFromConfig(cfg),
	}
}

func (c Client) ListVPCs(ownerId string) (VPCs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filters := []types.Filter{
		{Name: aws.String("state"), Values: []string{"available"}},
		{Name: aws.String("owner-id"), Values: []string{ownerId}},
	}
	in := &ec2.DescribeVpcsInput{Filters: filters}

	var vpcs VPCs
	for {
		out, err := c.svc.DescribeVpcs(ctx, in)
		if err != nil {
			return nil, err
		}
		vpcs = append(vpcs, toVPCs(out.Vpcs)...)
		if aws.ToString(out.NextToken) == "" {
			break
		}
		in.NextToken = out.NextToken
	}
	return vpcs, nil
}

func (c Client) CreateVPCFlowLogs(vpc VPC, logGroupName string, roleArn string, tags map[string]string) error {
	in := createFlowLogsInput{
		resourceType: types.FlowLogsResourceTypeVpc,
		resourceIds:  []string{vpc.Id},
		logGroupName: logGroupName,
		roleArn:      roleArn,
		tags:         tags,
	}
	return c.createFlowLogs(in)
}

func (c Client) ListSubnets(ownerId, vpcId string) (Subnets, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filters := []types.Filter{
		{Name: aws.String("state"), Values: []string{"available"}},
		{Name: aws.String("owner-id"), Values: []string{ownerId}},
		{Name: aws.String("vpc-id"), Values: []string{vpcId}},
	}
	in := &ec2.DescribeSubnetsInput{Filters: filters}

	var subnets Subnets
	for {
		out, err := c.svc.DescribeSubnets(ctx, in)
		if err != nil {
			return nil, err
		}
		subnets = append(subnets, toSubnets(out.Subnets)...)
		if aws.ToString(out.NextToken) == "" {
			break
		}
		in.NextToken = out.NextToken
	}
	return subnets, nil
}

func (c Client) CreateSubnetFlowLogs(subnet Subnet, logGroupName string, roleArn string, tags map[string]string) error {
	in := createFlowLogsInput{
		resourceType: types.FlowLogsResourceTypeSubnet,
		resourceIds:  []string{subnet.Id},
		logGroupName: logGroupName,
		roleArn:      roleArn,
		tags:         tags,
	}
	return c.createFlowLogs(in)
}

func (c Client) ListNatGateways(vpcId string) (NatGateways, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filters := []types.Filter{
		{Name: aws.String("vpc-id"), Values: []string{vpcId}},
	}

	in := &ec2.DescribeNatGatewaysInput{
		Filter: filters,
	}

	var natGateways NatGateways
	for {
		out, err := c.svc.DescribeNatGateways(ctx, in)
		if err != nil {
			return nil, err
		}
		natGateways = append(natGateways, toNatGateways(out.NatGateways)...)
		if aws.ToString(out.NextToken) == "" {
			break
		}
		in.NextToken = out.NextToken
	}
	return natGateways, nil
}

func (c Client) CreateNatGatewayFlowLogs(natGateway NatGateway, logGroupName string, roleArn string, tags map[string]string) error {
	in := createFlowLogsInput{
		resourceType: types.FlowLogsResourceTypeNetworkInterface,
		resourceIds:  []string{natGateway.NetworkInterfaceId},
		logGroupName: logGroupName,
		roleArn:      roleArn,
		tags:         tags,
	}
	return c.createFlowLogs(in)
}

func (c Client) ListSecurityGroups(ownerId, vpcId string) (SecurityGroups, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filters := []types.Filter{
		{Name: aws.String("owner-id"), Values: []string{ownerId}},
		{Name: aws.String("vpc-id"), Values: []string{vpcId}},
	}

	in := &ec2.DescribeSecurityGroupsInput{Filters: filters}

	var securityGroups SecurityGroups
	for {
		out, err := c.svc.DescribeSecurityGroups(ctx, in)
		if err != nil {
			return nil, err
		}
		securityGroups = append(securityGroups, toSecurityGroups(out.SecurityGroups)...)
		if aws.ToString(out.NextToken) == "" {
			break
		}
		in.NextToken = out.NextToken
	}
	return securityGroups, nil
}

func (c Client) CreateSecurityGroupFlowLogs(securityGroup SecurityGroup, logGroupName string, roleArn string, tags map[string]string) error {
	networkInterfaceIds, err := c.ListSecurityGroupNetworkInterfaceIds(securityGroup)
	if err != nil {
		return err
	}

	in := createFlowLogsInput{
		resourceType: types.FlowLogsResourceTypeNetworkInterface,
		resourceIds:  networkInterfaceIds,
		logGroupName: logGroupName,
		roleArn:      roleArn,
		tags:         tags,
	}
	return c.createFlowLogs(in)
}

func (c Client) ListInstances(vpcId string) (Instances, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filters := []types.Filter{
		{Name: aws.String("instance-state-name"), Values: []string{"running"}},
		{Name: aws.String("vpc-id"), Values: []string{vpcId}},
	}
	in := &ec2.DescribeInstancesInput{Filters: filters}

	var instances Instances
	for {
		out, err := c.svc.DescribeInstances(ctx, in)
		if err != nil {
			return nil, err
		}
		for _, reservation := range out.Reservations {
			instances = append(instances, toInstances(reservation.Instances)...)
		}
		if aws.ToString(out.NextToken) == "" {
			break
		}
		in.NextToken = out.NextToken
	}
	return instances, nil
}

func (c Client) CreateInstancesFlowLogs(instances Instances, logGroupName string, roleArn string, tags map[string]string) error {
	var networkInterfaceIds []string
	for _, v := range instances {
		networkInterfaceIds = append(networkInterfaceIds, v.NetworkInterfaceIds...)
	}

	in := createFlowLogsInput{
		resourceType: types.FlowLogsResourceTypeNetworkInterface,
		resourceIds:  networkInterfaceIds,
		logGroupName: logGroupName,
		roleArn:      roleArn,
		tags:         tags,
	}
	return c.createFlowLogs(in)
}

type createFlowLogsInput struct {
	resourceType types.FlowLogsResourceType
	resourceIds  []string
	logGroupName string
	roleArn      string
	tags         map[string]string
}

func (c createFlowLogsInput) toInput() *ec2.CreateFlowLogsInput {
	return &ec2.CreateFlowLogsInput{
		ResourceIds:              c.resourceIds,
		ResourceType:             c.resourceType,
		LogFormat:                aws.String(toLogFormat(fields.FlowLogFields)),
		LogGroupName:             aws.String(c.logGroupName),
		LogDestinationType:       types.LogDestinationTypeCloudWatchLogs,
		DeliverLogsPermissionArn: aws.String(c.roleArn),
		TrafficType:              types.TrafficTypeAll,
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeVpcFlowLog,
				Tags:         fromTags(c.tags),
			},
		},
	}
}

func (c Client) createFlowLogs(in createFlowLogsInput) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := c.svc.CreateFlowLogs(ctx, in.toInput()); err != nil {
		return err
	}
	return nil
}

// ListFlowLogs flow logs that match supplied tags and name (tag Name) prefix
func (c Client) ListFlowLogs(namePrefix string, tags map[string]string) (FlowLogs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// delete 'Name' from tags, we are using 'Name' as tag-key
	delete(tags, "Name")
	filters := []types.Filter{
		{Name: aws.String("log-destination-type"), Values: []string{"cloud-watch-logs"}},
		{Name: aws.String("tag-key"), Values: []string{"Name"}},
	}
	for k, v := range tags {
		filters = append(filters, types.Filter{Name: aws.String(fmt.Sprintf("tag:%s", k)), Values: []string{v}})
	}

	in := &ec2.DescribeFlowLogsInput{Filter: filters}

	var flowLogs FlowLogs
	for {
		out, err := c.svc.DescribeFlowLogs(ctx, in)
		if err != nil {
			return nil, err
		}
		flowLogs = append(flowLogs, toFlowLogs(out.FlowLogs)...)
		if aws.ToString(out.NextToken) == "" {
			break
		}
		in.NextToken = out.NextToken
	}

	if namePrefix == "" {
		c.logger.Debug("returning all flow logs, supplied prefix filter is empty")
		return flowLogs, nil
	}

	var filteredFlowLogs FlowLogs
	for _, flowLog := range flowLogs {
		name := flowLog.Name
		if strings.HasPrefix(name, namePrefix) {
			filteredFlowLogs = append(filteredFlowLogs, flowLog)
			continue
		}
		c.logger.Debug(fmt.Sprintf("flow log %s does not match %s name prefix", name, namePrefix))
	}
	return filteredFlowLogs, nil
}

func (c Client) DeleteFlowLogs(flowLogIds []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	in := &ec2.DeleteFlowLogsInput{FlowLogIds: flowLogIds}

	_, err := c.svc.DeleteFlowLogs(ctx, in)
	return err
}

func (c Client) ListSecurityGroupNetworkInterfaceIds(securityGroup SecurityGroup) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filters := []types.Filter{
		{Name: aws.String("vpc-id"), Values: []string{securityGroup.VpcId}},
		{Name: aws.String("group-id"), Values: []string{securityGroup.Id}},
	}

	in := &ec2.DescribeNetworkInterfacesInput{Filters: filters}

	var eniIds []string
	for {
		out, err := c.svc.DescribeNetworkInterfaces(ctx, in)
		if err != nil {
			return nil, fmt.Errorf("describe %s security group network interfaces: %w", securityGroup.Id, err)
		}
		for _, eni := range out.NetworkInterfaces {
			eniIds = append(eniIds, aws.ToString(eni.NetworkInterfaceId))
		}
		if aws.ToString(out.NextToken) == "" {
			break
		}
		in.NextToken = out.NextToken
	}
	return eniIds, nil
}

func fromTags(in map[string]string) []types.Tag {
	var out []types.Tag
	for k, v := range in {
		out = append(out, types.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	return out
}

func toLogFormat(logFields []string) string {
	var logFormat []string
	for _, v := range logFields {
		logFormat = append(logFormat, fmt.Sprintf("${%s}", v))
	}
	return strings.Join(logFormat, " ")
}