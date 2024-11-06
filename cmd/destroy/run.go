package destroy

import (
	"h3s/cmd/dependencies"
	"h3s/cmd/errors"

	"github.com/spf13/cobra"
)

func runDestroyCluster(_ *cobra.Command, _ []string) error {
	deps := dependencies.Get()

	ctx, err := deps.GetClusterContext()
	if err != nil {
		return errors.Wrap(errors.ErrorTypeCluster, "failed to load cluster context", err)
	}

	if err := deps.DestroyHetznerResources(ctx); err != nil {
		return errors.Wrap(errors.ErrorTypeHetzner, "failed to destroy cluster resources", err)
	}

	return nil
}
