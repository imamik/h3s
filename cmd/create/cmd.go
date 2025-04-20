// Package create provides commands for creating resources - a h3s cluster or necessary configuration & credentials
package create

import (
	"h3s/internal/utils/cli"

	"github.com/spf13/cobra"
)

// Command configurations
var (
	// Main create command configuration
	createConfig = cli.CommandConfig{
		Use:   "create",
		Short: "Create various resources",
		Long:  `Create various resources - Hetzner Cloud, k3s & configuration resources`,
	}

	// Config subcommand configuration
	configConfig = cli.CommandConfig{
		Use:   "config",
		Short: "Create a new cluster configuration",
		Long:  `Create a new cluster configuration. This command will prompt questions to configure the project and create a new configuration file.`,
		RunE:  runCreateConfig,
		Args:  cobra.NoArgs,
	}

	// Credentials subcommand configuration
	credentialsConfig = cli.CommandConfig{
		Use:   "credentials",
		Short: "Create cluster credentials",
		Long:  `Create cluster credentials. This command will prompt you various questions to configure the project and create a new credentials file.`,
		RunE:  runCreateCredentials,
		Args:  cobra.NoArgs,
	}

	// Cluster subcommand configuration
	clusterConfig = cli.CommandConfig{
		Use:   "cluster",
		Short: "Create a new cluster",
		Long:  `Create a new cluster - setup the necessary resources, install k3s and configure the cluster`,
		RunE:  runCreateCluster,
		Args:  cobra.NoArgs,
	}
)

// Cmd is the main command for creating resources
var Cmd *cobra.Command

// init initializes the create command and its subcommands
func init() {
	// Create subcommands
	configCmd := cli.NewCommand(configConfig)
	credentialsCmd := cli.NewCommand(credentialsConfig)
	clusterCmd := cli.NewCommand(clusterConfig)

	// Create main command with subcommands
	createConfig.Subcommands = []*cobra.Command{configCmd, credentialsCmd, clusterCmd}
	Cmd = cli.NewCommand(createConfig)
}
