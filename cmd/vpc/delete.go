package vpc

import (
	"fmt"
	"github.com/pete911/flowlogs/cmd/flag"
	"github.com/pete911/flowlogs/cmd/prompt"
	"github.com/pete911/flowlogs/internal/aws"
	"github.com/pete911/flowlogs/internal/aws/ec2"
	"github.com/spf13/cobra"
	"os"
)

var (
	Delete = &cobra.Command{
		Use:     "vpc",
		Aliases: []string{"vpcs"},
		Short:   "delete flow logs for specific vpc",
		Long:    "",
		Run:     runDelete,
	}
)

func runDelete(_ *cobra.Command, _ []string) {
	logger := flag.Global.Logger()
	client := aws.NewClient(logger, flag.Global.AWSConfig())

	selectedFlowLog := prompt.SelectFlowLog(prompt.ListFlowLogs(client, aws.FlowLogTypeVPC), true)
	if err := client.DeleteResources(ec2.FlowLogs{selectedFlowLog}); err != nil {
		fmt.Printf("delete flow logs: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("flow logs deleted")
}
