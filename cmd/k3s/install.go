package k3s

import (
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/pkg/clustercontext"
	"hcloud-k3s-cli/pkg/k3s/install"
	"hcloud-k3s-cli/pkg/resources/server"
)

var Install = &cobra.Command{
	Use:   "install",
	Short: "Install k3s on all servers in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := clustercontext.Context()
		servers := server.GetAll(ctx)
		for _, server := range servers {
			install.Install(ctx, server)
		}
	},
}
