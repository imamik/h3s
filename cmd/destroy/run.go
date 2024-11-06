package destroy

import (
	"h3s/internal/cluster"
	"h3s/internal/errors"
	"h3s/internal/hetzner"

	"github.com/spf13/cobra"
)

// runDestroyCluster destroys a h3s cluster (leaving the configuration files)
func runDestroyCluster(_ *cobra.Command, _ []string) error {
	// Load the cluster context
	ctx, err := cluster.Context()
	if err != nil {
		return errors.Wrap(errors.ErrorTypeCluster, "failed to load cluster context", err)
	}

	// Destroy the cluster resources
	if err := hetzner.Destroy(ctx); err != nil {
		return errors.Wrap(errors.ErrorTypeHetzner, "failed to destroy cluster resources", err)
	}

	return nil
}
