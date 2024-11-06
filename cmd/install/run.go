package install

import (
	"h3s/internal/cluster"
	"h3s/internal/errors"
	"h3s/internal/k3s"
	"h3s/internal/k8s"

	"github.com/spf13/cobra"
)

// runInstallK3s installs k3s on all servers in the h3s cluster
func runInstallK3s(_ *cobra.Command, _ []string) error {
	ctx, err := cluster.Context()
	if err != nil {
		return errors.Wrap(errors.ErrorTypeCluster, "failed to load cluster context", err)
	}
	if err := k3s.Install(ctx); err != nil {
		return errors.Wrap(errors.ErrorTypeK3s, "failed to install k3s", err)
	}
	return nil
}

// runInstallComponents installs all components on the h3s cluster
func runInstallComponents(_ *cobra.Command, _ []string) error {
	ctx, err := cluster.Context()
	if err != nil {
		return errors.Wrap(errors.ErrorTypeCluster, "failed to load cluster context", err)
	}
	if err := k8s.Install(ctx); err != nil {
		return errors.Wrap(errors.ErrorTypeK3s, "failed to install kubernetes components", err)
	}
	return nil
}
