package aws

import (
	"errors"
	"fmt"
	"github.com/pete911/flowlogs/internal/aws/ec2"
	"github.com/pete911/flowlogs/internal/aws/iam"
	"github.com/pete911/flowlogs/internal/aws/logs"
	"github.com/pete911/flowlogs/internal/aws/query"
	"log/slog"
)

type FlowLogType string

const (
	FlowLogTypeInstance      FlowLogType = "instance-"
	FlowLogTypeSecurityGroup FlowLogType = "sg-"
	FlowLogTypeNatGateway    FlowLogType = "nat-"
	FlowLogTypeSubnet        FlowLogType = "subnet-"
	FlowLogTypeVPC           FlowLogType = "vpc-"
	FlowLogTypeAll           FlowLogType = ""
)

type Client struct {
	config     Config
	logger     *slog.Logger
	ec2client  ec2.Client
	logsClient logs.Client
	iamClient  iam.Client
}

func NewClient(logger *slog.Logger, cfg Config) Client {
	return Client{
		config:     cfg,
		logger:     logger,
		ec2client:  ec2.NewClient(logger, cfg.Config),
		logsClient: logs.NewClient(logger, cfg.Config),
		iamClient:  iam.NewClient(logger, cfg.Config),
	}
}

func (c Client) ListFlowLogs(flowLogType FlowLogType) (ec2.FlowLogs, error) {
	// tag name does not matter, we are going to delete it and search by name prefix
	tags := NewTags("")
	delete(tags, "Name")
	// we use vpc id, which has 'vpc-' prefix
	return c.ec2client.ListFlowLogs(string(flowLogType), tags)
}

func (c Client) ListVPCs() (ec2.VPCs, error) {
	return c.ec2client.ListVPCs(c.config.Account)
}

func (c Client) CreateVPCFlowLogs(vpc ec2.VPC) (string, error) {
	tags := tagsFromId(vpc.Id)
	logGroupName, roleArn, err := c.createLogGroupAndRole(vpc.Id, tags)
	if err != nil {
		return "", err
	}

	if err := c.ec2client.CreateVPCFlowLogs(vpc, logGroupName, roleArn, tags); err != nil {
		return "", err
	}
	c.logger.Info("flow logs created")
	return logGroupName, err
}

func (c Client) ListSubnets(vpcId string) (ec2.Subnets, error) {
	return c.ec2client.ListSubnets(c.config.Account, vpcId)
}

func (c Client) CreateSubnetFlowLogs(subnet ec2.Subnet) (string, error) {
	tags := tagsFromId(subnet.Id)
	logGroupName, roleArn, err := c.createLogGroupAndRole(subnet.Id, tags)
	if err != nil {
		return "", err
	}

	if err := c.ec2client.CreateSubnetFlowLogs(subnet, logGroupName, roleArn, tags); err != nil {
		return "", err
	}
	c.logger.Info("flow logs created")
	return logGroupName, err
}

func (c Client) ListSecurityGroups(vpcId string) (ec2.SecurityGroups, error) {
	return c.ec2client.ListSecurityGroups(c.config.Account, vpcId)
}

func (c Client) CreateSecurityGroupFlowLogs(securityGroup ec2.SecurityGroup) (string, error) {
	tags := tagsFromId(securityGroup.Id)
	logGroupName, roleArn, err := c.createLogGroupAndRole(securityGroup.Id, tags)
	if err != nil {
		return "", err
	}

	if err := c.ec2client.CreateSecurityGroupFlowLogs(securityGroup, logGroupName, roleArn, tags); err != nil {
		return "", err
	}
	c.logger.Info("flow logs created")
	return logGroupName, err
}

func (c Client) ListNatGateways(vpcId string) (ec2.NatGateways, error) {
	return c.ec2client.ListNatGateways(vpcId)
}

func (c Client) CreateNatGatewayFlowLogs(natGateway ec2.NatGateway) (string, error) {
	tags := tagsFromId(natGateway.Id)
	logGroupName, roleArn, err := c.createLogGroupAndRole(natGateway.Id, tags)
	if err != nil {
		return "", err
	}

	if err := c.ec2client.CreateNatGatewayFlowLogs(natGateway, logGroupName, roleArn, tags); err != nil {
		return "", err
	}
	c.logger.Info("flow logs created")
	return logGroupName, err
}

func (c Client) ListInstances(vpcId string) (ec2.Instances, error) {
	return c.ec2client.ListInstances(vpcId)
}

func (c Client) CreateInstanceFlowLogs(instances ec2.Instances) (string, error) {
	if err := validateInstances(instances); err != nil {
		return "", err
	}

	// we don't use id only for instances, because they are grouped by name, instead we prefix name with 'instance-'
	id := fmt.Sprintf("%s%s", FlowLogTypeInstance, instances[0].Name)
	tags := tagsFromId(id)
	logGroupName, roleArn, err := c.createLogGroupAndRole(id, tags)
	if err != nil {
		return "", err
	}

	if err := c.ec2client.CreateInstancesFlowLogs(instances, logGroupName, roleArn, tags); err != nil {
		return "", err
	}
	c.logger.Info("flow logs created")
	return logGroupName, err
}

// createLogGroupAndRole creates cloud watch log group and IAM role and returns log group and IAM role ARN
func (c Client) createLogGroupAndRole(id string, tags map[string]string) (string, string, error) {
	logGroupName := logGroupNameFromId(id)
	if err := c.logsClient.CreateLogGroup(logGroupName, tags); err != nil {
		return "", "", err
	}
	c.logger.Info(fmt.Sprintf("log group %s created", logGroupName))

	roleName := c.iamRoleNameFromId(id)
	roleArn, err := c.iamClient.CreateFlowLogsRole(roleName, tags)
	if err != nil {
		return "", "", err
	}
	c.logger.Info(fmt.Sprintf("role %s created", roleName))
	return logGroupName, roleArn, err
}

// QueryFlowLogs run query on specified flow logs
func (c Client) QueryFlowLogs(flowLogs ec2.FlowLogs, query query.Query) ([]map[string]string, error) {
	if len(flowLogs) == 0 {
		c.logger.Info("no flow logs provided, nothing to query")
		return nil, nil
	}

	var logGroupNames []string
	for _, v := range flowLogs {
		logGroupNames = append(logGroupNames, logGroupNameFromFlowLogName(v.Name))
	}
	return c.logsClient.Query(logGroupNames, query.GetQuery(), query.GetSinceMinutes(), query.GetLimit())
}

func (c Client) ListNetworkInterfaces() (ec2.NetworkInterfaces, error) {
	return c.ec2client.ListNetworkInterfaces()
}

// DeleteResources delete flow logs, IAM roles and cloud watch log groups
func (c Client) DeleteResources(flowLogs ec2.FlowLogs) error {
	if len(flowLogs) == 0 {
		c.logger.Info("no flow logs provided, nothing to delete")
		return nil
	}

	if err := c.DeleteFlowLogs(flowLogs); err != nil {
		return err
	}
	if err := c.DeleteIAMRoles(flowLogs); err != nil {
		return err
	}
	if err := c.DeleteLogGroups(flowLogs); err != nil {
		return err
	}
	return nil
}

func (c Client) DeleteFlowLogs(flowLogs ec2.FlowLogs) error {
	if len(flowLogs) == 0 {
		c.logger.Info("no flow logs provided, nothing to delete")
		return nil
	}

	if err := c.ec2client.DeleteFlowLogs(flowLogs.Ids()); err != nil {
		return fmt.Errorf("delete flow logs: %w", err)
	}
	c.logger.Info(fmt.Sprintf("%d flow logs deleted", len(flowLogs)))
	return nil
}

func (c Client) DeleteIAMRoles(flowLogs ec2.FlowLogs) error {
	if len(flowLogs) == 0 {
		c.logger.Info("no flow logs provided, nothing to delete")
		return nil
	}

	names, flowLogsByName := flowLogs.GetByNames()
	c.logger.Debug(fmt.Sprintf("deleting %d roles", len(names)))
	for _, name := range names {
		if len(flowLogsByName[name]) == 0 {
			continue
		}
		// get tags, but without name
		tags := flowLogsByName[name][0].Tags()
		delete(tags, "Name")
		roleName := c.iamRoleNameFromFlowLogName(name)
		if err := c.iamClient.DeleteRole(roleName, tags); err != nil {
			return fmt.Errorf("delete iam role: %w", err)
		}
		c.logger.Info(fmt.Sprintf("%s iam role deleted", roleName))
	}
	return nil
}

func (c Client) DeleteLogGroups(flowLogs ec2.FlowLogs) error {
	if len(flowLogs) == 0 {
		c.logger.Info("no flow logs provided, nothing to delete")
		return nil
	}

	names, flowLogsByName := flowLogs.GetByNames()
	c.logger.Debug(fmt.Sprintf("deleting %d log groups", len(names)))
	for _, name := range names {
		if len(flowLogsByName[name]) == 0 {
			continue
		}
		// get tags, but without name
		tags := flowLogsByName[name][0].Tags()
		delete(tags, "Name")
		logGroupName := logGroupNameFromFlowLogName(name)
		if err := c.logsClient.DeleteLogGroup(logGroupName, tags); err != nil {
			return fmt.Errorf("delete log group: %w", err)
		}
		c.logger.Info(fmt.Sprintf("%s log group deleted", logGroupName))
	}
	return nil
}

func (c Client) iamRoleNameFromId(id string) string {
	// IAM is global, we need to include region in the name
	roleName := fmt.Sprintf("fl-cli-%s-%s", c.config.Region, id)
	return trimRoleName(roleName)
}

func (c Client) iamRoleNameFromFlowLogName(flowLogName string) string {
	// prefix e.g. 'instance-', 'vpc-' etc. is already included in the flow log name
	roleName := fmt.Sprintf("fl-cli-%s-%s", c.config.Region, flowLogName)
	return trimRoleName(roleName)
}

// trimRoleName trim if longer than 64 chars - max role name is 64 characters
func trimRoleName(roleName string) string {
	if len(roleName) > 64 {
		return roleName[:64]
	}
	return roleName
}

func logGroupNameFromId(id string) string {
	return fmt.Sprintf("/fl-cli/%s", id)
}

func logGroupNameFromFlowLogName(flowLogName string) string {
	return fmt.Sprintf("/fl-cli/%s", flowLogName)
}

func tagsFromId(id string) map[string]string {
	return NewTags(id)
}

// validateInstances checks if instances are not 0 size and they are the same name
func validateInstances(instances ec2.Instances) error {
	if len(instances) == 0 {
		return errors.New("no instances provided")
	}
	for i := range instances {
		if i == 0 {
			continue
		}
		// make sure every instance in this 'group' has the same name
		one, two := instances[i-1].Name, instances[i].Name
		if one != two {
			return fmt.Errorf("supplied instances do not have the same name: %s %s", one, two)
		}
	}
	return nil
}
