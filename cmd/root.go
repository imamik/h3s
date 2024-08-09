// Package cmd contains the command-line interface definitions for the h3s tool
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"h3s/cmd/cluster"
	"h3s/cmd/config"
	"h3s/cmd/credentials"
	"h3s/cmd/k3s"
	"h3s/cmd/kubectl"
	"h3s/cmd/ssh"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "h3s",
	Short: "A CLI to setup k3s Kubernetes resources on Hetzner Cloud",
	Long:  "h3s (Hetzner Highly-Available-k3s Clusters) is a command-line interface for setting up and managing k3s Kubernetes resources on Hetzner Cloud. It provides various subcommands for managing clusters, configurations, and resources.",
	Run:   printWelcome,
}

// printWelcome prints a welcome message & help information when the root command is called
func printWelcome(_ *cobra.Command, _ []string) {
	fmt.Println("Welcome to h3s CLI")
	fmt.Println("Use --help for more information about available commands")
}

// init function sets up the command structure by adding subcommands to the root command
func init() {
	RootCmd.AddCommand(
		config.Cmd,
		credentials.Cmd,
		cluster.Cmd,
		k3s.Cmd,
		kubectl.Cmd,
		ssh.Cmd,
	)
}
