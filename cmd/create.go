package cmd

import (
	"github.com/pete911/flowlogs/cmd/instance"
	"github.com/pete911/flowlogs/cmd/nat"
	"github.com/pete911/flowlogs/cmd/sg"
	"github.com/pete911/flowlogs/cmd/subnet"
	"github.com/pete911/flowlogs/cmd/vpc"
	"github.com/spf13/cobra"
)

var (
	Create = &cobra.Command{
		Use:   "create",
		Short: "create AWS flow logs",
		Long:  "",
	}
)

func init() {
	Root.AddCommand(Create)
	Create.AddCommand(instance.Create)
	Create.AddCommand(nat.Create)
	Create.AddCommand(sg.Create)
	Create.AddCommand(subnet.Create)
	Create.AddCommand(vpc.Create)
}
