package credentials

import (
	"github.com/spf13/cobra"
)

// Cmd is the command to configure project credentials (Hetzner Cloud & DNS API token & k3s token)
var Cmd = &cobra.Command{
	Use:   "credentials",
	Short: "Manage project credentials",
}

// createCommand is the command to configure project credentials (Hetzner Cloud & DNS API token & k3s token)
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create project credentials",
	RunE:  runCreate,
}

// init adds subcommands to the Cluster command
func init() {
	Cmd.AddCommand(createCmd)
}
