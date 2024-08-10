package cluster

import (
	"github.com/spf13/cobra"
	"h3s/internal/clustercontext"
	"h3s/internal/k3s/install"
	"h3s/internal/resources/cluster"
)

func runCreate(_ *cobra.Command, _ []string) error {
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

// runDestroy is the function that is executed when the destroy command is called - it destroys the cluster resources
func runDestroy(_ *cobra.Command, _ []string) error {
	// Get the cluster context
	ctx := clustercontext.Context()
	// Destroy the cluster resources
	cluster.Destroy(ctx)
	return nil
}
