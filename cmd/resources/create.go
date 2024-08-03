package resources

import (
	"github.com/spf13/cobra"
	"h3s/internal/clustercontext"
	"h3s/internal/resources/cluster"
)

var Create = &cobra.Command{
	Use:   "create",
	Short: "Create a new resources",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := clustercontext.Context()
		cluster.Create(ctx)
	},
}
