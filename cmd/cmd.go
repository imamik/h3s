// Package cmd contains the command-line interface definitions for the h3s tool
package cmd

import (
	"github.com/spf13/cobra"
	"h3s/cmd/cluster"
	"h3s/cmd/config"
	"h3s/cmd/credentials"
	"h3s/cmd/k3s"
	"h3s/cmd/kubectl"
	"h3s/cmd/ssh"
	"h3s/internal/version"
)

var (
	Version = version.Create(0, 1, 0)
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "h3s",
	Short:   "A CLI to setup k3s Kubernetes resources on Hetzner Cloud",
	Long:    "h3s (Hetzner Highly-Available-k3s Clusters) is a command-line interface for setting up and managing k3s Kubernetes resources on Hetzner Cloud. It provides various subcommands for managing clusters, configurations, and resources.",
	Version: Version.String(),
	RunE:    runRoot,
}

// init function sets up the command structure with all high level subcommands & sets up the flags
func init() {
	// Add subcommands
	RootCmd.AddCommand(
		cluster.Cmd,
		config.Cmd,
		credentials.Cmd,
		k3s.Cmd,
		kubectl.Cmd,
		ssh.Cmd,
	)
}
