package prompt

import (
	"fmt"
	"github.com/pete911/flowlogs/internal/aws"
	"github.com/pete911/flowlogs/internal/aws/ec2"
	"os"
)

func ListVPCs(client aws.Client) ec2.VPCs {
	vpcs, err := client.ListVPCs()
	if err != nil {
		fmt.Printf("list vpcs: %v\n", err)
		os.Exit(1)
	}

	existingVPCFlowLogs, err := client.ListFlowLogs(aws.FlowLogTypeVPC)
	if err != nil {
		fmt.Printf("list vpc flow logs: %v\n", err)
		os.Exit(1)
	}
	return vpcs.FilterOut(existingVPCFlowLogs)
}

func SelectVPC(vpcs ec2.VPCs, confirm bool) ec2.VPC {
	if len(vpcs) == 0 {
		fmt.Println("no vpcs found")
		os.Exit(1)
	}

	var items []string
	for _, vpc := range vpcs {
		items = append(items, vpc.String())
	}
	label := fmt.Sprintf("select vpc [%d vpcs]:", len(vpcs))
	i, _ := Select(label, items)
	if confirm {
		Confirm("selected vpc, continue")
	}
	return vpcs[i]
}
