package destroy

import (
	"github.com/spf13/cobra"
	"h3s/internal/cluster"
	"h3s/internal/hetzner"
)

// runDestroyCluster destroys a h3s cluster (leaving the configuration files)
func runDestroyCluster(_ *cobra.Command, _ []string) error {
	// Load the cluster context
	ctx, err := cluster.Context()
	if err != nil {
		return err
	}

	// Destroy the cluster resources
	err = hetzner.Destroy(ctx)
	if err != nil {
		return err
	}

	return nil
}
