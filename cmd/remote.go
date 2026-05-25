package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// bit remote add <name> bit://<network>/<registry>/<repo-id>
var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Manage remotes",
}

var remoteAddCmd = &cobra.Command{
	Use:   "add <name> <url>",
	Short: "Add a new remote",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, url := args[0], args[1]
		fmt.Printf("TODO: save remote %s -> %s\n", name, url)
		return nil
	},
}

func init() {
	remoteCmd.AddCommand(remoteAddCmd)
}
