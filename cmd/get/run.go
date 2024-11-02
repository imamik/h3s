package get

import (
	"fmt"
	"h3s/internal/cluster"
	"h3s/internal/k8s"
	"h3s/internal/k8s/kubeconfig"
	"h3s/internal/utils/execute"

	"github.com/spf13/cobra"
)

// runGetKubeConfig gets the kubeconfig for the h3s cluster
func runGetKubeConfig(_ *cobra.Command, _ []string) error {
	ctx, err := cluster.Context()
	if err != nil {
		return err
	}
	return kubeconfig.Download(ctx)
}

// runGetToken gets a fresh bearer token for the h3s cluster
func runGetToken(cmd *cobra.Command, _ []string) error {
	ctx, err := cluster.Context()
	if err != nil {
		return err
	}

	// Create a new bearer token for the k3s dashboard
	b, err := k8s.Token(ctx, "kubernetes-dashboard", "admin-user", 24)
	if err != nil {
		cmd.PrintErrf("Failed to get bearer token: %v\n", err)
		return err
	}

	// Build the command to copy the bearer token to the clipboard
	localCmd := fmt.Sprintf("printf '%%s' \"%s\" | pbcopy", b)
	if _, err := execute.Local(localCmd); err != nil {
		cmd.PrintErrf("Failed to copy bearer token to clipboard: %v\n", err)
		return err
	}

	cmd.Println("Bearer token copied to clipboard.")
	return nil
}
