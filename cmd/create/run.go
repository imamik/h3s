package create

import (
	"github.com/spf13/cobra"
	"h3s/internal/cluster"
	"h3s/internal/config/create"
	"h3s/internal/config/credentials"
	"h3s/internal/hetzner"
	"h3s/internal/k3s"
	"h3s/internal/k3s/install"
)

// runCreateConfig creates a new h3s cluster configuration
func runCreateConfig(cmd *cobra.Command, _ []string) error {
	// Load the latest 5 stable k3s releases
	k3sReleases, err := k3s.GetFilteredReleases(false, true, 5)

	// If there was an error, print it and return
	if err != nil {
		return err
	}

	// Build the configuration
	create.Build(k3sReleases)

	return nil
}

// runCreateCredentials creates h3s cluster credentials
func runCreateCredentials(_ *cobra.Command, _ []string) error {
	credentials.Configure()
	return nil
}

// runCreateCluster creates a h3s cluster
func runCreateCluster(_ *cobra.Command, _ []string) error {
	// Load the cluster context
	ctx, err := cluster.Context()
	if err != nil {
		return err
	}

	// Create the cluster resources
	hetzner.Create(ctx)

	// Install k3s on the cluster
	install.K3s(ctx)

	// Install additional software on the cluster (traefik, cert-manager, csi etc.)
	install.Software(ctx)

	// Download the kubeconfig file
	install.DownloadKubeconfig(ctx)

	return nil
}
