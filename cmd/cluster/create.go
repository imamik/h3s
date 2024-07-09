package cluster

import (
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/k3s/install"
	"hcloud-k3s-cli/internal/resources/cluster"
)

var k3sInstall bool
var cleanup bool

var Create = &cobra.Command{
	Use:   "create",
	Short: "Create a new resources",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := clustercontext.Context()
		cluster.Create(ctx)
		if k3sInstall {
			install.Install(ctx, cleanup)
		}
	},
}

func init() {
	Create.Flags().BoolVar(&cleanup, "cleanup", true, "Force installation")
	Create.Flags().BoolVar(&k3sInstall, "install", true, "Install k3s on all servers in the cluster")
}
