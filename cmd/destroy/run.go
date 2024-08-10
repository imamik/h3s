package destroy

import (
	"github.com/spf13/cobra"
	"h3s/internal/clustercontext"
	"h3s/internal/resources/cluster"
)

// runDestroyCluster destroys a h3s cluster (leaving the configuration files)
func runDestroyCluster(_ *cobra.Command, _ []string) error {
	// Get the cluster context
	ctx := clustercontext.Context()
	// Destroy the cluster resources
	cluster.Destroy(ctx)
	return nil
}
