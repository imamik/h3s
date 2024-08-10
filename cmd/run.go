package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// runRoot prints information - version, if flag is set or welcome message with help info when called without any arguments
func runRoot(_ *cobra.Command, _ []string) error {
	fmt.Println("Welcome to h3s CLI")
	fmt.Println("Use --help for more information about available commands")
	return nil
}
