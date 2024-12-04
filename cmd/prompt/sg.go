package prompt

import (
	"fmt"
	"github.com/pete911/flowlogs/internal/aws"
	"github.com/pete911/flowlogs/internal/aws/ec2"
	"os"
)

func ListSecurityGroups(client aws.Client, vpcId string) ec2.SecurityGroups {
	sgs, err := client.ListSecurityGroups(vpcId)
	if err != nil {
		fmt.Printf("list security groups: %v\n", err)
		os.Exit(1)
	}

	existingSGFlowLogs, err := client.ListFlowLogs(aws.FlowLogTypeSecurityGroup)
	if err != nil {
		fmt.Printf("list security group flow logs: %v\n", err)
		os.Exit(1)
	}
	return sgs.FilterOut(existingSGFlowLogs)
}

func SelectSecurityGroup(sgs ec2.SecurityGroups, confirm bool) ec2.SecurityGroup {
	if len(sgs) == 0 {
		fmt.Println("no security group found")
		os.Exit(1)
	}

	var items []string
	for _, sg := range sgs {
		items = append(items, sg.String())
	}
	label := fmt.Sprintf("select security gruops [%d security groups]:", len(sgs))
	i, _ := Select(label, items)
	if confirm {
		Confirm("selected security group, continue")
	}
	return sgs[i]
}
