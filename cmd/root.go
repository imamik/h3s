package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"h3s/cmd/cluster"
	"h3s/cmd/config"
	"h3s/cmd/k3s"
	"h3s/cmd/kubectl"
	"h3s/cmd/microos"
	"h3s/cmd/resources"
	"h3s/cmd/ssh"
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
	rootCmd.AddCommand(ssh.Ssh)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
