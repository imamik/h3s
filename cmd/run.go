package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// runRoot prints information - version, if flag is set or welcome message with help info when called without any arguments
func runRoot(cmd *cobra.Command, _ []string) error {
	_, err := fmt.Fprintln(cmd.OutOrStdout(), "Welcome to h3s CLI\nUse --help for more information about available commands")
	if err != nil {
		return err
	}
	return nil
}
