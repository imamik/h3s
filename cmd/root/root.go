// Package root provides the root command and initialization for the h3s CLI
package root

import (
	"fmt"

	"github.com/spf13/cobra"

	"h3s/internal/version"
)

// Cmd is the root command for the h3s CLI
var Cmd *cobra.Command

// runRoot prints information - version, if flag is set or welcome message with help info when called without any arguments
func runRoot(cmd *cobra.Command, _ []string) {
	cmd.Println("Welcome to h3s CLI")
	cmd.Println("Use --help for more information about available commands")
}

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
}
