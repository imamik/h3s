package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/cmd/cfg"
	"hcloud-k3s-cli/cmd/cluster"
	"hcloud-k3s-cli/cmd/k3s"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "hcloud-k8s",
	Short: "A CLI to setup a k3s Kubernetes cluster on Hetzner Cloud",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to my CLI")
	},
}

func init() {
	rootCmd.AddCommand(k3s.K3s)
	rootCmd.AddCommand(cfg.Cfg)
	rootCmd.AddCommand(cluster.Cluster)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
