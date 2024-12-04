package cmd

import (
	"fmt"
	"github.com/pete911/flowlogs/cmd/flag"
	"github.com/pete911/flowlogs/internal/aws"
	"github.com/spf13/cobra"
	"os"
)

var (
	list = &cobra.Command{
		Use:   "list",
		Short: "list all resources created by cli",
		Long:  "",
		Run:   runList,
	}
)

func init() {
	Root.AddCommand(list)
}

func runList(_ *cobra.Command, _ []string) {
	logger := flag.Global.Logger()
	client := aws.NewClient(logger, flag.Global.AWSConfig())

	flowLogs, err := client.ListFlowLogs(aws.FlowLogTypeAll)
	if err != nil {
		fmt.Printf("list flow logs: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%d flow logs found\n", len(flowLogs))
	for _, flowLog := range flowLogs {
		fmt.Println(flowLog)
	}
}
