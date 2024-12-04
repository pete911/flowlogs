package prompt

import (
	"fmt"
	"github.com/pete911/flowlogs/internal/aws"
	"github.com/pete911/flowlogs/internal/aws/ec2"
	"os"
)

func ListSubnets(client aws.Client, vpcId string) ec2.Subnets {
	subnets, err := client.ListSubnets(vpcId)
	if err != nil {
		fmt.Printf("list subnets: %v\n", err)
		os.Exit(1)
	}

	existingSubnetFlowLogs, err := client.ListFlowLogs(aws.FlowLogTypeSubnet)
	if err != nil {
		fmt.Printf("list subnet flow logs: %v\n", err)
		os.Exit(1)
	}
	return subnets.FilterOut(existingSubnetFlowLogs)
}

func SelectSubnet(subnets ec2.Subnets, confirm bool) ec2.Subnet {
	if len(subnets) == 0 {
		fmt.Println("no subnet found")
		os.Exit(1)
	}

	var items []string
	for _, subnet := range subnets {
		items = append(items, subnet.String())
	}
	label := fmt.Sprintf("select subnet [%d subnets]:", len(subnets))
	i, _ := Select(label, items)
	if confirm {
		Confirm("selected subnet, continue")
	}
	return subnets[i]
}
