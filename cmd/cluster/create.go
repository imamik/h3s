package cluster

import (
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/pkg/clustercontext"
	"hcloud-k3s-cli/pkg/k3s/install"
	"hcloud-k3s-cli/pkg/resources/cluster"
)

var Create = &cobra.Command{
	Use:   "create",
	Short: "Create a new resources",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := clustercontext.Context()
		cluster.Create(ctx)
		install.Install(ctx)
	},
}
