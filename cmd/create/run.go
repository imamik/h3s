package create

import (
	"github.com/spf13/cobra"
	"h3s/internal/cluster"
	"h3s/internal/config/create"
	"h3s/internal/config/credentials"
	"h3s/internal/hetzner"
	"h3s/internal/k3s"
	"h3s/internal/k8s"
	"h3s/internal/k8s/kubeconfig"
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
	clr, err := cluster.Context()
	if err != nil {
		return err
	}

	// Create the cluster resources
	hetzner.Create(clr)

	// Install k3s on the cluster
	if err := k3s.Install(clr); err != nil {
		return err
	}

	// Install additional software on the cluster (traefik, cert-manager, csi etc.)
	if err := k8s.Install(clr); err != nil {
		return err
	}

	// Download the kubeconfig file
	if err := kubeconfig.Download(clr); err != nil {
		return err
	}

	return nil
}