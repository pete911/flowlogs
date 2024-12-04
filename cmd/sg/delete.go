package sg

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
		Use:     "sg",
		Aliases: []string{"sgs", "security-group", "security-groups"},
		Short:   "delete flow logs for specific security group",
		Long:    "",
		Run:     runDelete,
	}
)

func runDelete(_ *cobra.Command, _ []string) {
	logger := flag.Global.Logger()
	client := aws.NewClient(logger, flag.Global.AWSConfig())

	selectedFlowLogs := prompt.SelectFlowLogs(prompt.ListFlowLogs(client, aws.FlowLogTypeSecurityGroup), true)
	if err := client.DeleteResources(selectedFlowLogs); err != nil {
		fmt.Printf("delete flow logs: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("flow logs deleted")
}
