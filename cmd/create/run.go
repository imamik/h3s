package create

import (
	"h3s/cmd/dependencies"
	"h3s/cmd/errors"

	"github.com/spf13/cobra"
)

func runCreateConfig(_ *cobra.Command, _ []string) error {
	deps := dependencies.Get()

	k3sReleases, err := deps.GetK3sReleases(false, true, 5)
	if err != nil {
		return errors.Wrap(errors.ErrorTypeK3s, "failed to get k3s releases", err)
	}

	if err := deps.BuildClusterConfig(k3sReleases); err != nil {
		return errors.Wrap(errors.ErrorTypeConfig, "failed to build configuration", err)
	}

	return nil
}

func runCreateCredentials(_ *cobra.Command, _ []string) error {
	deps := dependencies.Get()

	if err := deps.ConfigureCredentials(); err != nil {
		return errors.Wrap(errors.ErrorTypeConfig, "failed to configure credentials", err)
	}
	return nil
}

func runCreateCluster(_ *cobra.Command, _ []string) error {
	deps := dependencies.Get()

	ctx, err := deps.GetClusterContext()
	if err != nil {
		return errors.Wrap(errors.ErrorTypeCluster, "failed to load cluster context", err)
	}

	if err := deps.CreateHetznerResources(ctx); err != nil {
		return errors.Wrap(errors.ErrorTypeHetzner, "failed to create hetzner resources", err)
	}

	if err := deps.InstallK3s(ctx); err != nil {
		return errors.Wrap(errors.ErrorTypeK3s, "failed to install k3s", err)
	}

	if err := deps.InstallK8sComponents(ctx); err != nil {
		return errors.Wrap(errors.ErrorTypeK3s, "failed to install kubernetes components", err)
	}

	if err := deps.DownloadKubeconfig(ctx); err != nil {
		return errors.Wrap(errors.ErrorTypeK3s, "failed to download kubeconfig file", err)
	}

	return nil
}
