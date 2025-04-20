// Package root provides the root command and initialization for the h3s CLI
package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cmdversion "h3s/cmd/version"
)

// Cmd is the root command for the h3s CLI
var Cmd *cobra.Command

// runRoot prints information - version, if flag is set or welcome message with help info when called without any arguments
func runRoot(cmd *cobra.Command, _ []string) {
	cmd.Println("Welcome to h3s CLI")
	cmd.Println("Use --help for more information about available commands")
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	Cmd = &cobra.Command{
		Use:   "h3s",
		Short: "A CLI to setup k3s Kubernetes resources on Hetzner Cloud",
		Long: `h3s simplifies the creation, management, and deletion of K3s clusters
hosted on Hetzner Cloud infrastructure.`,
		Run: runRoot,
	}

	Cmd.AddCommand(cmdversion.Cmd)

	if err := Cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
