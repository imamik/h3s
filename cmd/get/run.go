package get

import (
	"fmt"
	"github.com/spf13/cobra"
	"h3s/internal/clustercontext"
	"h3s/internal/k3s/install"
	"h3s/internal/k3s/kubeconfig"
	"h3s/internal/k3s/token"
	"os/exec"
)

// runGetKubeConfig gets the kubeconfig for the h3s cluster
func runGetKubeConfig(_ *cobra.Command, _ []string) error {
	ctx := clustercontext.Context()
	_, _, gatewayServer, controlPlaneNodes, _ := install.GetSetup(ctx)
	kubeconfig.Download(ctx, gatewayServer, controlPlaneNodes[0])
	return nil
}

// runGetToken gets a fresh bearer token for the h3s cluster
func runGetToken(_ *cobra.Command, _ []string) error {
	ctx := clustercontext.Context()
	b, err := token.Create(ctx, "kubernetes-dashboard", "admin-user", 24)
	if err != nil {
		fmt.Printf("Failed to get bearer token: %v\n", err)
		return err
	}

	// Correctly handling the bearer token with special characters
	copyCmd := exec.Command("sh", "-c", fmt.Sprintf("printf '%%s' \"%s\" | pbcopy", b))
	if err := copyCmd.Run(); err != nil {
		fmt.Printf("Failed to copy bearer token to clipboard: %v\n", err)
		return err
	}

	fmt.Println("Bearer token copied to clipboard.")
	return nil
}
