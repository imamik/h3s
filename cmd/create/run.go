package create

import (
	"github.com/spf13/cobra"
	"h3s/internal/clustercontext"
	"h3s/internal/config/create"
	"h3s/internal/config/credentials"
	"h3s/internal/k3s/install"
	"h3s/internal/k3s/releases"
	"h3s/internal/resources/cluster"
)

// runCreateConfig creates a new h3s cluster configuration
func runCreateConfig(cmd *cobra.Command, _ []string) error {
	// Get the latest 5 stable k3s releases
	k3sReleases, err := releases.GetFilteredReleases(false, true, 5)

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
	// Get the cluster context
	ctx := clustercontext.Context()

	// Create the cluster resources
	cluster.Create(ctx)

	// Install k3s on the cluster
	install.K3s(ctx)

	// Install additional software on the cluster (traefik, cert-manager, csi etc.)
	install.Software(ctx)

	// Download the kubeconfig file
	install.DownloadKubeconfig(ctx)

	return nil
}
