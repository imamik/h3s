package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/cmd/config"
	"hcloud-k3s-cli/cmd/k3s"
	"hcloud-k3s-cli/pkg/client"
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

var metadataCmd = &cobra.Command{
	Use:   "k3s",
	Short: "Commands related to k3s",
}

func init() {
	client.InitHcloudClient()
	rootCmd.AddCommand(k3s.ReleasesCmd)
	rootCmd.AddCommand(config.InitCmd)

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
