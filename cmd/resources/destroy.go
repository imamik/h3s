package resources

import (
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/cluster"
)

var Destroy = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy a new resources",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := clustercontext.Context()
		cluster.Destroy(ctx)
	},
}
