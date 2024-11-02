package install

import (
	"h3s/internal/cluster"
	"h3s/internal/k3s"
	"h3s/internal/k8s"

	"github.com/spf13/cobra"
)

// runInstallK3s installs k3s on all servers in the h3s cluster
func runInstallK3s(_ *cobra.Command, _ []string) error {
	ctx, err := cluster.Context()
	if err != nil {
		return err
	}
	return k3s.Install(ctx)
}

// runInstallComponents installs all components on the h3s cluster
func runInstallComponents(_ *cobra.Command, _ []string) error {
	ctx, err := cluster.Context()
	if err != nil {
		return err
	}
	return k8s.Install(ctx)
}
