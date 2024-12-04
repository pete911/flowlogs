package vpc

import (
	"fmt"
	"github.com/pete911/flowlogs/cmd/flag"
	"github.com/pete911/flowlogs/cmd/out"
	"github.com/pete911/flowlogs/cmd/prompt"
	"github.com/pete911/flowlogs/internal/aws"
	"github.com/pete911/flowlogs/internal/aws/fields"
	"github.com/spf13/cobra"
	"os"
)

var (
	Query = &cobra.Command{
		Use:     "vpc",
		Aliases: []string{"vpcs"},
		Short:   "query flow logs for vpc",
		Long:    "",
		Run:     runQuery,
	}
)

func runQuery(_ *cobra.Command, _ []string) {
	logger := flag.Global.Logger()
	client := aws.NewClient(logger, flag.Global.AWSConfig())

	selectedFlowLogs := prompt.SelectFlowLogs(prompt.ListFlowLogs(client, aws.FlowLogTypeVPC), false)
	logs, err := client.QueryFlowLogs(selectedFlowLogs, flag.Query.GetQuery())
	if err != nil {
		fmt.Printf("query flow logs: %v\n", err)
		os.Exit(1)
	}

	tableRows := fields.QueryFields.Values(logs)
	table := out.NewTable(logger, os.Stdout)
	table.AddRow(fields.QueryFields.Header()...)
	for _, row := range tableRows {
		table.AddRow(row...)
	}
	table.Print()
}