// Package cluster provides commands for managing k3s Kubernetes clusters on Hetzner Cloud
package cluster

import (
	"github.com/spf13/cobra"
)

// Cmd is the main command for cluster management
var Cmd = &cobra.Command{
	Use:   "cluster",
	Short: "Manage k3s Kubernetes clusters on Hetzner Cloud",
	Long:  `Manage k3s Kubernetes clusters on Hetzner Cloud - create, destroy, list and get information about clusters`,
}

// Create is a cobra.Command that handles the creation of a new cluster
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new cluster",
	Long:  `Create a new cluster - setup the necessary resources, install k3s and configure the cluster`,
	RunE:  runCreate,
}

// Destroy is a cobra.Command that handles the destruction of an existing cluster
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy an existing cluster",
	Long:  `Destroy an existing cluster including alle resources`,
	RunE:  runDestroy,
}

// init adds subcommands to the Cluster command
func init() {
	Cmd.AddCommand(createCmd, destroyCmd)
}
