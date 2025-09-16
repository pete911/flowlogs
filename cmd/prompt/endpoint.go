package prompt

import (
	"fmt"
	"os"

	"github.com/pete911/flowlogs/internal/aws"
	"github.com/pete911/flowlogs/internal/aws/ec2"
)

func ListVPCEndpoints(client aws.Client) ec2.VPCEndpoints {
	endpoints, err := client.ListVPCEndpoints()
	if err != nil {
		fmt.Printf("list vpc endpoints: %v\n", err)
		os.Exit(1)
	}

	existingVPCFlowLogs, err := client.ListFlowLogs(aws.FlowLogTypeVPCEndpoint)
	if err != nil {
		fmt.Printf("list vpc endpoint flow logs: %v\n", err)
		os.Exit(1)
	}
	return endpoints.FilterOut(existingVPCFlowLogs)
}

func SelectVPCEndpoint(endpoints ec2.VPCEndpoints, confirm bool) ec2.VPCEndpoint {
	if len(endpoints) == 0 {
		fmt.Println("no vpc endpoints found")
		os.Exit(1)
	}

	var items []string
	for _, endpoint := range endpoints {
		items = append(items, endpoint.String())
	}
	label := fmt.Sprintf("select vpc endpoint [%d endpoints]:", len(endpoints))
	i, _ := Select(label, items)
	if confirm {
		Confirm("selected vpc endpoint, continue")
	}
	return endpoints[i]
}
