package install

import (
	"h3s/cmd/dependencies"
	"h3s/cmd/errors"

	"github.com/spf13/cobra"
)

func runInstallK3s(_ *cobra.Command, _ []string) error {
	deps := dependencies.Get()

	ctx, err := deps.GetClusterContext()
	if err != nil {
		return errors.Wrap(errors.ErrorTypeCluster, "failed to load cluster context", err)
	}

	if err := deps.InstallK3s(ctx); err != nil {
		return errors.Wrap(errors.ErrorTypeK3s, "failed to install k3s", err)
	}
	return nil
}

func runInstallComponents(_ *cobra.Command, _ []string) error {
	deps := dependencies.Get()

	ctx, err := deps.GetClusterContext()
	if err != nil {
		return errors.Wrap(errors.ErrorTypeCluster, "failed to load cluster context", err)
	}

	if err := deps.InstallK8sComponents(ctx); err != nil {
		return errors.Wrap(errors.ErrorTypeK3s, "failed to install kubernetes components", err)
	}
	return nil
}
