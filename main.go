// main.go
package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "hcloud-k8s-cli",
	Short: "A CLI to setup a k3s Kubernetes cluster on Hetzner Cloud",
	Long:  `A CLI to setup a k3s Kubernetes cluster on Hetzner Cloud`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to my CLI")
	},
}

func main() {
	initClient()
	rootCmd.AddCommand(listCmd) // Register the list command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
