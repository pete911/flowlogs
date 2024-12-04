package prompt

import (
	"fmt"
	"github.com/pete911/flowlogs/internal/aws"
	"github.com/pete911/flowlogs/internal/aws/ec2"
	"os"
)

func ListNatGateways(client aws.Client, vpcId string) ec2.NatGateways {
	sgs, err := client.ListNatGateways(vpcId)
	if err != nil {
		fmt.Printf("list nat gateways: %v\n", err)
		os.Exit(1)
	}

	existingNATFlowLogs, err := client.ListFlowLogs(aws.FlowLogTypeNatGateway)
	if err != nil {
		fmt.Printf("list nat gateway flow logs: %v\n", err)
		os.Exit(1)
	}
	return sgs.FilterOut(existingNATFlowLogs)
}

func SelectNatGateway(natGateways ec2.NatGateways, confirm bool) ec2.NatGateway {
	if len(natGateways) == 0 {
		fmt.Println("no nat gateways found")
		os.Exit(1)
	}

	var items []string
	for _, nat := range natGateways {
		items = append(items, nat.String())
	}
	label := fmt.Sprintf("select nat gateways [%d nat gateway]:", len(natGateways))
	i, _ := Select(label, items)
	if confirm {
		Confirm("selected nat gateway, continue")
	}
	return natGateways[i]
}
