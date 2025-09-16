package endpoint

import (
	"fmt"
	"os"

	"github.com/pete911/flowlogs/cmd/flag"
	"github.com/pete911/flowlogs/cmd/prompt"
	"github.com/pete911/flowlogs/internal/aws"
	"github.com/spf13/cobra"
)

var (
	Delete = &cobra.Command{
		Use:     "endpoint",
		Aliases: []string{"endpoints", "vpc-endpoint", "vpc-endpoints"},
		Short:   "delete flow logs for specific endpoint",
		Long:    "",
		Run:     runDelete,
	}
)

func runDelete(_ *cobra.Command, _ []string) {
	logger := flag.Global.Logger()
	client := aws.NewClient(logger, flag.Global.AWSConfig())

	selectedFlowLogs := prompt.SelectFlowLogs(prompt.ListFlowLogs(client, aws.FlowLogTypeVPCEndpoint), true)
	if err := client.DeleteResources(selectedFlowLogs); err != nil {
		fmt.Printf("delete flow logs: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("flow logs deleted")
}
