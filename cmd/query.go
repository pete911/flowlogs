package cmd

import (
	"github.com/pete911/flowlogs/cmd/flag"
	"github.com/pete911/flowlogs/cmd/instance"
	"github.com/pete911/flowlogs/cmd/nat"
	"github.com/pete911/flowlogs/cmd/sg"
	"github.com/pete911/flowlogs/cmd/subnet"
	"github.com/pete911/flowlogs/cmd/vpc"
	"github.com/spf13/cobra"
)

var (
	Query = &cobra.Command{
		Use:   "query",
		Short: "query AWS flow logs",
		Long:  "",
	}
)

func init() {
	flag.InitPersistentQueryFlags(Query, &flag.Query)
	Root.AddCommand(Query)
	Query.AddCommand(instance.Query)
	Query.AddCommand(nat.Query)
	Query.AddCommand(sg.Query)
	Query.AddCommand(subnet.Query)
	Query.AddCommand(vpc.Query)
}
