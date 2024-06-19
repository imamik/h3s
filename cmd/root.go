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
	Use:   "hcloud-k8s-cli",
	Short: "A CLI to setup a k3s Kubernetes cluster on Hetzner Cloud",
	Long:  `A CLI to setup a k3s Kubernetes cluster on Hetzner Cloud`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to my CLI")
	},
}

func init() {
	rootCmd.AddCommand(k3s.ReleasesCmd)
	rootCmd.AddCommand(cfg.CreateCmd)
	rootCmd.AddCommand(cluster.ClusterCmd)

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
