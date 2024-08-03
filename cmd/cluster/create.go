package cluster

import (
	"github.com/spf13/cobra"
	"h3s/internal/clustercontext"
	"h3s/internal/k3s/install"
	"h3s/internal/resources/cluster"
)

var k3sInstall bool

var Create = &cobra.Command{
	Use:   "create",
	Short: "Create a new resources",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := clustercontext.Context()
		cluster.Create(ctx)
		if k3sInstall {
			install.Install(ctx)
		}
	},
}

func init() {
	Create.Flags().BoolVar(&k3sInstall, "install", true, "Install k3s on all servers in the cluster")
}
