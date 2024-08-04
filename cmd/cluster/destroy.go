package cluster

import (
	"github.com/spf13/cobra"
	"h3s/internal/clustercontext"
	"h3s/internal/resources/cluster"
)

// Destroy is a cobra.Command that handles the destruction of an existing cluster
var Destroy = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy an existing cluster",
	Long:  `Destroy an existing cluster including alle resources`,
	Run:   runDestroy,
}

func runDestroy(cmd *cobra.Command, args []string) {
	// Get the cluster context
	ctx := clustercontext.Context()
	// Destroy the cluster resources
	cluster.Destroy(ctx)
}
