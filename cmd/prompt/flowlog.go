package prompt

import (
	"fmt"
	"github.com/pete911/flowlogs/internal/aws"
	"github.com/pete911/flowlogs/internal/aws/ec2"
	"os"
)

func ListFlowLogs(client aws.Client, flowLogType aws.FlowLogType) ec2.FlowLogs {
	flowLogs, err := client.ListFlowLogs(flowLogType)
	if err != nil {
		fmt.Printf("delete flow logs: %v\n", err)
		os.Exit(1)
	}
	return flowLogs
}

func SelectFlowLog(flowLogs ec2.FlowLogs, confirm bool) ec2.FlowLog {
	if len(flowLogs) == 0 {
		fmt.Println("no flow logs found")
		os.Exit(1)
	}

	var items []string
	for _, flowLog := range flowLogs {
		items = append(items, flowLog.String())
	}
	i, _ := Select("select flow log:", items)
	if confirm {
		Confirm("selected flow log, continue")
	}
	return flowLogs[i]
}

func SelectFlowLogs(flowLogs ec2.FlowLogs, confirm bool) ec2.FlowLogs {
	if len(flowLogs) == 0 {
		fmt.Println("no flow logs found")
		os.Exit(1)
	}

	names, flowLogsByName := flowLogs.GetByNames()
	var items []string
	for _, name := range names {
		items = append(items, fmt.Sprintf("%s [%d flow logs]", name, len(flowLogsByName[name])))
	}
	label := fmt.Sprintf("select flow logs [%d groups]:", len(flowLogsByName))
	i, _ := Select(label, items)

	selectedFlowLogs := flowLogsByName[names[i]]
	fmt.Println("selected flow logs:")
	for _, selectedFlowLog := range selectedFlowLogs {
		fmt.Println(selectedFlowLog)
	}
	if confirm {
		Confirm("selected flow logs, continue")
	}
	return selectedFlowLogs
}
