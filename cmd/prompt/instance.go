package prompt

import (
	"fmt"
	"github.com/pete911/flowlogs/internal/aws"
	"github.com/pete911/flowlogs/internal/aws/ec2"
	"os"
)

func ListInstances(client aws.Client, vpcId string) ec2.Instances {
	instances, err := client.ListInstances(vpcId)
	if err != nil {
		fmt.Printf("list instances: %v\n", err)
		os.Exit(1)
	}

	existingInstanceFlowLogs, err := client.ListFlowLogs(aws.FlowLogTypeInstance)
	if err != nil {
		fmt.Printf("list instance flow logs: %v\n", err)
		os.Exit(1)
	}
	return instances.FilterOut(existingInstanceFlowLogs)
}

func SelectInstances(instances ec2.Instances, confirm bool) ec2.Instances {
	if len(instances) == 0 {
		fmt.Println("no instances found")
		os.Exit(1)
	}

	names, instancesByName := instances.GetByNames()
	var items []string
	for _, name := range names {
		items = append(items, fmt.Sprintf("%s [%d instances]", name, len(instancesByName[name])))
	}
	label := fmt.Sprintf("select instances [%d groups]:", len(instancesByName))
	i, _ := Select(label, items)

	selectedInstances := instancesByName[names[i]]
	fmt.Println("selected instances:")
	for _, selectedInstance := range selectedInstances {
		fmt.Println(selectedInstance)
	}
	if confirm {
		Confirm("selected instances, continue")
	}
	return selectedInstances
}
