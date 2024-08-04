// Package cluster provides commands for managing k3s Kubernetes clusters on Hetzner Cloud
package cluster

import (
	"github.com/spf13/cobra"
)

// Cluster is the main command for cluster management
var Cluster = &cobra.Command{
	Use:   "cluster",
	Short: "Manage k3s Kubernetes clusters on Hetzner Cloud",
	Long:  `Manage k3s Kubernetes clusters on Hetzner Cloud - create, destroy, list and get information about clusters`,
}

// init function is automatically called when the package is initialized
// It adds subcommands to the Cluster command
func init() {
	// Add Create subcommand to Cluster command
	Cluster.AddCommand(Create)
	// Add Destroy subcommand to Cluster command
	Cluster.AddCommand(Destroy)
}
