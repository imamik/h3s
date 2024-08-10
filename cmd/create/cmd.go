package create

import (
	"github.com/spf13/cobra"
)

// Cmd is the main command for creating resources - Hetzner Cloud, k3s & configuration resources
var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Create various resources",
	Long:  `Create various resources - Hetzner Cloud, k3s & configuration resources`,
}

// createConfigCmd creates a project configuration for a h3s cluster
var createConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Create a new cluster configuration",
	Long:  `Create a new cluster configuration. This command will prompt questions to configure the project and create a new configuration file.`,
	RunE:  runCreateConfig,
}

// createCredentialsCmd creates project credentials (Hetzner Cloud & DNS API token & k3s token) for a h3s cluster
var createCredentialsCmd = &cobra.Command{
	Use:   "credentials",
	Short: "Create cluster credentials",
	Long:  `Create cluster credentials. This command will prompt you various questions to configure the project and create a new credentials file.`,
	RunE:  runCreateCredentials,
}

// createClusterCmd creates a new h3s cluster
var createClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Create a new cluster",
	Long:  `Create a new cluster - setup the necessary resources, install k3s and configure the cluster`,
	RunE:  runCreateCluster,
}

// init adds subcommands to the create command
func init() {
	Cmd.AddCommand(
		createConfigCmd,
		createCredentialsCmd,
		createClusterCmd,
	)
}
