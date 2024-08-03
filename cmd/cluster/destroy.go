package cluster

import (
	"github.com/spf13/cobra"
	"h3s/internal/clustercontext"
	"h3s/internal/resources/cluster"
)

var Destroy = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy a new resources",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := clustercontext.Context()
		cluster.Destroy(ctx)
	},
}
