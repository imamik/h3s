// Package cmd contains the command-line interface definitions for the h3s tool
package cmd

import (
	"fmt"
	"h3s/cmd/create"
	"h3s/cmd/destroy"
	"h3s/cmd/get"
	"h3s/cmd/install"
	"h3s/cmd/kubectl"
	"h3s/cmd/ssh"
	"h3s/internal/version"

	versioncmd "h3s/cmd/version"

	"github.com/spf13/cobra"
)

// Cmd is the root command for the h3s CLI
var Cmd *cobra.Command

// Initialize sets up the root command with version
func Initialize(info version.BuildInfo) {
	Cmd = &cobra.Command{
		Use: "h3s",
		Version: fmt.Sprintf("%s\nCommit: %s\nGo version: %s",
			info.Version,
			info.Commit,
			info.GoVersion,
		),
		Short: "A CLI to setup k3s Kubernetes resources on Hetzner Cloud",
		Long:  "h3s (Hetzner Highly-Available-k3s Clusters) is a command-line interface for setting up and managing k3s Kubernetes resources on Hetzner Cloud. It provides various subcommands for managing clusters, configurations, and resources.",
		Run:   runRoot,
	}

	// Add version flags
	Cmd.Flags().BoolP("version", "v", false, "Print version information")

	// Add subcommands
	Cmd.AddCommand(
		versioncmd.Cmd,
		create.Cmd,
		destroy.Cmd,
		get.Cmd,
		install.Cmd,
		kubectl.Cmd,
		ssh.Cmd,
	)
}
