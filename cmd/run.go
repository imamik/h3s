package cmd

import (
	"github.com/spf13/cobra"
)

// runRoot prints information - version, if flag is set or welcome message with help info when called without any arguments
func runRoot(cmd *cobra.Command, _ []string) {
	cmd.Println("Welcome to h3s CLI")
	cmd.Println("Use --help for more information about available commands")
}
