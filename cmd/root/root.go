// Package root provides the root command and initialization for the h3s CLI
package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cmdcreate "h3s/cmd/create"
	cmddestroy "h3s/cmd/destroy"
	cmdget "h3s/cmd/get"
	cmdinstall "h3s/cmd/install"
	cmdkubectl "h3s/cmd/kubectl"
	cmdssh "h3s/cmd/ssh"
	cmdversion "h3s/cmd/version"
	internalversion "h3s/internal/version"
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
	var showVersion bool
	Cmd = &cobra.Command{
		Use:   "h3s",
		Short: "A CLI to setup k3s Kubernetes resources on Hetzner Cloud",
		Long: `h3s simplifies the creation, management, and deletion of K3s clusters
hosted on Hetzner Cloud infrastructure.`,
		Run: func(cmd *cobra.Command, args []string) {
			if showVersion {
				info := internalversion.GetBuildInfo()
				cmd.Printf("h3s version %s\n", info.Version)
				cmd.Printf("Commit: %s\n", info.Commit)
				cmd.Printf("Go version: %s\n", info.GoVersion)
				return
			}
			runRoot(cmd, args)
		},
	}

	Cmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "Show version information and exit")

	Cmd.AddCommand(cmdversion.Cmd)
	Cmd.AddCommand(cmdcreate.Cmd)
	Cmd.AddCommand(cmddestroy.Cmd)
	Cmd.AddCommand(cmdget.Cmd)
	Cmd.AddCommand(cmdinstall.Cmd)
	Cmd.AddCommand(cmdkubectl.Cmd)
	Cmd.AddCommand(cmdssh.Cmd)

	if err := Cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
