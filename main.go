// main.go
package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "mycli",
	Short: "A brief description of your application",
	Long:  `A longer description of your application with usage examples.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to my CLI")
	},
}

func main() {
	rootCmd.AddCommand(listCmd) // Register the list command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
