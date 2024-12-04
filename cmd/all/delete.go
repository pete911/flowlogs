package all

import (
	"fmt"
	"github.com/pete911/flowlogs/cmd/flag"
	"github.com/pete911/flowlogs/cmd/prompt"
	"github.com/pete911/flowlogs/internal/aws"
	"github.com/spf13/cobra"
	"os"
)

var (
	Delete = &cobra.Command{
		Use:   "all",
		Short: "delete all resources created by cli",
		Long:  "",
		Run:   runDelete,
	}
)

func runDelete(_ *cobra.Command, _ []string) {
	logger := flag.Global.Logger()
	client := aws.NewClient(logger, flag.Global.AWSConfig())

	flowLogs := prompt.ListFlowLogs(client, aws.FlowLogTypeAll)
	if len(flowLogs) == 0 {
		fmt.Println("no flow logs found, nothing to delete")
		return
	}

	fmt.Printf("%d flow logs found\n", len(flowLogs))
	for _, flowLog := range flowLogs {
		fmt.Println(flowLog)
	}

	prompt.Confirm("selected flow logs, continue")

	if err := client.DeleteResources(flowLogs); err != nil {
		fmt.Printf("delete flow logs: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("flow logs deleted")
}
