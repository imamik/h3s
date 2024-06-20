package cluster

import (
	"github.com/spf13/cobra"
)

var Cluster = &cobra.Command{
	Use:   "cluster",
	Short: "CLI to manage k3s Kubernetes clusters on Hetzner Cloud",
}

func init() {
	Cluster.AddCommand(Create)
	Cluster.AddCommand(Destroy)
}
