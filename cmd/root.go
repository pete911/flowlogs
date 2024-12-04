package cmd

import (
	"github.com/pete911/flowlogs/cmd/flag"
	"github.com/spf13/cobra"
)

var (
	Root = &cobra.Command{}

	Version string
)

func init() {
	flag.InitPersistentFlags(Root, &flag.Global)
}
