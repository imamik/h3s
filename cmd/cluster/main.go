package cluster

import (
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/pkg/client"
)

var ClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "CLI to manage k3s Kubernetes clusters on Hetzner Cloud",
}

func init() {
	client.InitHcloudClient()
	ClusterCmd.AddCommand(CreateCmd)
}
