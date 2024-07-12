package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/cmd/cluster"
	"hcloud-k3s-cli/cmd/config"
	"hcloud-k3s-cli/cmd/k3s"
	"hcloud-k3s-cli/cmd/kubectl"
	"hcloud-k3s-cli/cmd/microos"
	"hcloud-k3s-cli/cmd/resources"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "hcloud-k8s",
	Short: "A CLI to setup a k3s Kubernetes resources on Hetzner Cloud",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to my CLI")
	},
}

func init() {
	rootCmd.AddCommand(k3s.K3s)
	rootCmd.AddCommand(config.Config)
	rootCmd.AddCommand(cluster.Cluster)
	rootCmd.AddCommand(microos.Image)
	rootCmd.AddCommand(resources.Resources)
	rootCmd.AddCommand(kubectl.Kubectl)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
