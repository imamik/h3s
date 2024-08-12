package get

import (
	"fmt"
	"github.com/spf13/cobra"
	"h3s/internal/cluster"
	"h3s/internal/k3s/install"
	"h3s/internal/k3s/kubeconfig"
	"h3s/internal/k3s/token"
	"os/exec"
)

// runGetKubeConfig gets the kubeconfig for the h3s cluster
func runGetKubeConfig(_ *cobra.Command, _ []string) error {
	ctx, err := cluster.Context()
	if err != nil {
		return err
	}
	_, _, gatewayServer, controlPlaneNodes, _ := install.GetSetup(ctx)
	kubeconfig.Download(ctx, gatewayServer, controlPlaneNodes[0])
	return nil
}

// runGetToken gets a fresh bearer token for the h3s cluster
func runGetToken(cmd *cobra.Command, _ []string) error {
	ctx, err := cluster.Context()
	if err != nil {
		return err
	}

	// Create a new bearer token for the k3s dashboard
	b, err := token.Create(ctx, "kubernetes-dashboard", "admin-user", 24)
	if err != nil {
		cmd.PrintErrf("Failed to get bearer token: %v\n", err)
		return err
	}

	// Build the command to copy the bearer token to the clipboard
	copyCmd := exec.Command("sh", "-c", fmt.Sprintf("printf '%%s' \"%s\" | pbcopy", b))

	// Run the command
	if err := copyCmd.Run(); err != nil {
		cmd.PrintErrf("Failed to copy bearer token to clipboard: %v\n", err)
		return err
	}

	cmd.Println("Bearer token copied to clipboard.")
	return nil
}
