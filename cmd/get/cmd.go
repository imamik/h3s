// Package get contains the commands for getting information from the cluster e.g. kubeconfig, token, etc.
package get

import (
	"github.com/spf13/cobra"
)

// Cmd is the main command for getting information from the cluster
var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Load information about the cluster",
	Long:  `Load various information about the cluster - kubeconfig, token, and other information`,
}

// getKubeConfigCmd gets the kubeconfig for the h3s cluster
var getKubeConfigCmd = &cobra.Command{
	Use:   "kubeconfig",
	Short: "Load the kubeconfig for the k3s cluster",
	RunE:  runGetKubeConfig,
}

// getTokenCmd gets the bearer token for the k3s dashboard
var getTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Load the bearer token for the k3s cluster",
	RunE:  runGetToken,
}

// init adds subcommands to the Cluster command
func init() {
	Cmd.AddCommand(
		getTokenCmd,
		getKubeConfigCmd,
	)
}
