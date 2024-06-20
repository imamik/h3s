package cluster

import (
	"github.com/spf13/cobra"
)

var ClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "CLI to manage k3s Kubernetes clusters on Hetzner Cloud",
}

func init() {
	ClusterCmd.AddCommand(CreateCmd)
}
