package config

import (
	"github.com/spf13/cobra"
)

// Cmd is the command to manage project configuration
var Cmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Project configuration",
	Long:  `Manage Project configuration - for now only there is only the "create" subcommand`,
}

// createCmd is the command to configure project configuration for a h3s cluster
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new h3s cluster configuration",
	Long:  `Create a new h3s cluster configuration. This command will prompt you various questions to configure the project and create a new configuration file.`,
	RunE:  runCreate,
}

// init adds subcommands to the Cluster command
func init() {
	Cmd.AddCommand(createCmd)
}
