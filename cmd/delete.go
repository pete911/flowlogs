package cmd

import (
	"github.com/pete911/flowlogs/cmd/all"
	"github.com/pete911/flowlogs/cmd/endpoint"
	"github.com/pete911/flowlogs/cmd/instance"
	"github.com/pete911/flowlogs/cmd/nat"
	"github.com/pete911/flowlogs/cmd/sg"
	"github.com/pete911/flowlogs/cmd/subnet"
	"github.com/pete911/flowlogs/cmd/vpc"
	"github.com/spf13/cobra"
)

var (
	Delete = &cobra.Command{
		Use:   "delete",
		Short: "delete AWS flow logs",
		Long:  "",
	}
)

func init() {
	Root.AddCommand(Delete)
	Delete.AddCommand(all.Delete)
	Delete.AddCommand(instance.Delete)
	Delete.AddCommand(nat.Delete)
	Delete.AddCommand(sg.Delete)
	Delete.AddCommand(subnet.Delete)
	Delete.AddCommand(vpc.Delete)
	Delete.AddCommand(endpoint.Delete)
}
