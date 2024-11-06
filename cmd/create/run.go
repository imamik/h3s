package create

import (
	"h3s/internal/cluster"
	"h3s/internal/config/create"
	"h3s/internal/config/credentials"
	"h3s/internal/errors"
	"h3s/internal/hetzner"
	"h3s/internal/k3s"
	"h3s/internal/k8s"
	"h3s/internal/k8s/kubeconfig"

	"github.com/spf13/cobra"
)

// runCreateConfig creates a new h3s cluster configuration
func runCreateConfig(cmd *cobra.Command, _ []string) error {
	// Load the latest 5 stable k3s releases
	k3sReleases, err := k3s.GetFilteredReleases(false, true, 5)
	if err != nil {
		return errors.Wrap(errors.ErrorTypeK3s, "failed to get k3s releases", err)
	}

	if err := create.Build(k3sReleases); err != nil {
		return errors.Wrap(errors.ErrorTypeConfig, "failed to build configuration", err)
	}

	return nil
}

// runCreateCredentials creates h3s cluster credentials
func runCreateCredentials(_ *cobra.Command, _ []string) error {
	return credentials.Configure()
}

// runCreateCluster creates a h3s cluster
func runCreateCluster(_ *cobra.Command, _ []string) error {
	// Load the cluster context
	clr, err := cluster.Context()
	if err != nil {
		return err
	}

	// Create the cluster resources
	if err := hetzner.Create(clr); err != nil {
		return errors.Wrap(errors.ErrorTypeHetzner, "failed to create hetzner resources", err)
	}

	// Install k3s on the cluster
	if err := k3s.Install(clr); err != nil {
		return errors.Wrap(errors.ErrorTypeK3s, "failed to install k3s", err)
	}

	// Install additional software on the cluster (traefik, cert-manager, csi etc.)
	if err := k8s.Install(clr); err != nil {
		return errors.Wrap(errors.ErrorTypeK3s, "failed to install kubernetes components", err)
	}

	// Download the kubeconfig file
	if err := kubeconfig.Download(clr); err != nil {
		return errors.Wrap(errors.ErrorTypeK3s, "failed to download kubeconfig file", err)
	}

	return nil
}
