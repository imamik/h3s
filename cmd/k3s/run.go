package k3s

import (
	"fmt"
	"github.com/spf13/cobra"
	"h3s/internal/clustercontext"
	"h3s/internal/k3s/bearer"
	"h3s/internal/k3s/install"
	"h3s/internal/k3s/kubeconfig"
	"h3s/internal/k3s/releases"
	"os/exec"
)

func runReleases(_ *cobra.Command, _ []string) error {
	r, err := releases.GetFilteredReleases(prerelease, stable, limit)
	if err != nil {
		println("Error fetching releases:", err)
		return err
	}
	releases.PrintReleases(r)
	return nil
}

func runKubeConfig(_ *cobra.Command, _ []string) error {
	ctx := clustercontext.Context()
	_, _, gatewayServer, controlPlaneNodes, _ := install.GetSetup(ctx)
	kubeconfig.Download(ctx, gatewayServer, controlPlaneNodes[0])
	return nil
}

func runInstall(_ *cobra.Command, _ []string) error {
	ctx := clustercontext.Context()
	install.K3s(ctx)
	install.Software(ctx)
	install.DownloadKubeconfig(ctx)
	return nil
}

func runInstallSoftware(_ *cobra.Command, _ []string) error {
	ctx := clustercontext.Context()
	install.Software(ctx)
	return nil
}

func runToken(_ *cobra.Command, _ []string) error {
	ctx := clustercontext.Context()
	b, err := bearer.GetBearerToken(ctx, "kubernetes-dashboard", "admin-user", 24)
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
