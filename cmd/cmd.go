// Package cmd contains the command-line interface definitions for the h3s tool
package cmd

import (
	"h3s/cmd/create"
	"h3s/cmd/destroy"
	"h3s/cmd/get"
	"h3s/cmd/install"
	"h3s/cmd/kubectl"
	"h3s/cmd/ssh"
	"h3s/internal/utils/version"

	"github.com/spf13/cobra"
)

// Version represents the current version of the h3s CLI tool
var Version = version.New(0, 1, 0).String()

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:     "h3s",
	Short:   "A CLI to setup k3s Kubernetes resources on Hetzner Cloud",
	Long:    "h3s (Hetzner Highly-Available-k3s Clusters) is a command-line interface for setting up and managing k3s Kubernetes resources on Hetzner Cloud. It provides various subcommands for managing clusters, configurations, and resources.",
	Version: Version,
	Run:     runRoot,
}

// init function sets up the command structure with all high level subcommands & sets up the flags
func init() {
	// Add subcommands
	Cmd.AddCommand(
		create.Cmd,
		destroy.Cmd,
		get.Cmd,
		install.Cmd,
		kubectl.Cmd,
		ssh.Cmd,
	)
}
