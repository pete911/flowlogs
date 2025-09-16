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
	Create = &cobra.Command{
		Use:     "endpoint",
		Aliases: []string{"endpoints", "vpc-endpoint", "vpc-endpoints"},
		Short:   "create flow logs for specific vpc endpoint",
		Long:    "",
		Run:     runCreate,
	}
)

func runCreate(_ *cobra.Command, _ []string) {
	logger := flag.Global.Logger()
	client := aws.NewClient(logger, flag.Global.AWSConfig())

	selectedEndpoint := prompt.SelectVPCEndpoint(prompt.ListVPCEndpoints(client), false)
	logGroup, err := client.CreateVPCEndpointFlowLogs(selectedEndpoint)
	if err != nil {
		fmt.Printf("create flow logs: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("created %s log group\n", logGroup)
}
