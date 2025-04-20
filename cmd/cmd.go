// Package cmd contains the command-line interface definitions for the h3s tool
package cmd

import (
	"h3s/cmd/create"
	"h3s/cmd/destroy"
	"h3s/cmd/get"
	"h3s/cmd/install"
	"h3s/cmd/kubectl"
	"h3s/cmd/root"
	"h3s/cmd/ssh"
	"h3s/internal/version"

	"github.com/spf13/cobra"

	versioncmd "h3s/cmd/version"
)

// Cmd is the root command for the h3s CLI
var Cmd *cobra.Command

// Initialize sets up the root command with version
func Initialize(_ version.BuildInfo) {
	root.Execute() // Execute handles command setup now
	Cmd = root.Cmd

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
