package instance

import (
	"fmt"
	"github.com/pete911/flowlogs/cmd/flag"
	"github.com/pete911/flowlogs/cmd/prompt"
	"github.com/pete911/flowlogs/internal/aws"
	"github.com/spf13/cobra"
	"os"
)

var (
	Create = &cobra.Command{
		Use:     "instance",
		Aliases: []string{"instances", "ec2", "ec2s"},
		Short:   "create flow logs for specific instance",
		Long:    "",
		Run:     runCreate,
	}
)

func runCreate(_ *cobra.Command, _ []string) {
	logger := flag.Global.Logger()
	client := aws.NewClient(logger, flag.Global.AWSConfig())

	selectedVPC := prompt.SelectVPC(prompt.ListVPCs(client), false)
	selectedInstances := prompt.SelectInstances(prompt.ListInstances(client, selectedVPC.Id), true)
	logGroup, err := client.CreateInstanceFlowLogs(selectedInstances)
	if err != nil {
		fmt.Printf("create flow logs: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("created %s log group\n", logGroup)
}
