package install

import (
	"github.com/spf13/cobra"
	"h3s/internal/cluster"
	"h3s/internal/k3s/install"
)

// runInstallK3s installs k3s on all servers in the h3s cluster
func runInstallK3s(_ *cobra.Command, _ []string) error {
	ctx, err := cluster.Context()
	if err != nil {
		return err
	}

	install.K3s(ctx)
	install.Software(ctx)
	install.DownloadKubeconfig(ctx)
	return nil
}

// runInstallComponents installs all components on the h3s cluster
func runInstallComponents(_ *cobra.Command, _ []string) error {
	ctx, err := cluster.Context()
	if err != nil {
		return err
	}
	install.Software(ctx)
	return nil
}
