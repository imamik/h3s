package get

import (
	"fmt"
	"h3s/internal/cluster"
	"h3s/internal/errors"
	"h3s/internal/k8s"
	"h3s/internal/k8s/kubeconfig"
	"h3s/internal/utils/execute"

	"github.com/spf13/cobra"
)

// runGetKubeConfig gets the kubeconfig for the h3s cluster
func runGetKubeConfig(_ *cobra.Command, _ []string) error {
	ctx, err := cluster.Context()
	if err != nil {
		return errors.Wrap(errors.ErrorTypeCluster, "failed to load cluster context", err)
	}
	if err := kubeconfig.Download(ctx); err != nil {
		return errors.Wrap(errors.ErrorTypeK3s, "failed to download kubeconfig", err)
	}
	return nil
}

// runGetToken gets a fresh bearer token for the h3s cluster
func runGetToken(cmd *cobra.Command, _ []string) error {
	ctx, err := cluster.Context()
	if err != nil {
		return errors.Wrap(errors.ErrorTypeCluster, "failed to load cluster context", err)
	}

	// Create a new bearer token for the k3s dashboard
	b, err := k8s.Token(ctx, "kubernetes-dashboard", "admin-user", 24)
	if err != nil {
		return errors.Wrap(errors.ErrorTypeK3s, "failed to generate bearer token", err)
	}

	// Build the command to copy the bearer token to the clipboard
	localCmd := fmt.Sprintf("printf '%%s' \"%s\" | pbcopy", b)
	if _, err := execute.Local(localCmd); err != nil {
		return errors.Wrap(errors.ErrorTypeSystem, "failed to copy bearer token to clipboard", err)
	}

	cmd.Println("Bearer token copied to clipboard.")
	return nil
}
