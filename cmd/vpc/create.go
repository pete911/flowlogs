package vpc

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
		Use:     "vpc",
		Aliases: []string{"vpcs"},
		Short:   "create flow logs for specific vpc",
		Long:    "",
		Run:     runCreate,
	}
)

func runCreate(_ *cobra.Command, _ []string) {
	logger := flag.Global.Logger()
	client := aws.NewClient(logger, flag.Global.AWSConfig())

	selectedVPC := prompt.SelectVPC(prompt.ListVPCs(client), true)
	logGroup, err := client.CreateVPCFlowLogs(selectedVPC)
	if err != nil {
		fmt.Printf("create flow logs: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("created %s log group\n", logGroup)
}
