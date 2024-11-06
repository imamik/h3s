package get

import (
	"fmt"
	"h3s/cmd/dependencies"
	"h3s/cmd/errors"

	"github.com/spf13/cobra"
)

func runGetKubeConfig(_ *cobra.Command, _ []string) error {
	deps := dependencies.Get()

	ctx, err := deps.GetClusterContext()
	if err != nil {
		return errors.Wrap(errors.ErrorTypeCluster, "failed to load cluster context", err)
	}

	if err := deps.DownloadKubeconfig(ctx); err != nil {
		return errors.Wrap(errors.ErrorTypeK3s, "failed to download kubeconfig", err)
	}
	return nil
}

func runGetToken(cmd *cobra.Command, _ []string) error {
	deps := dependencies.Get()

	ctx, err := deps.GetClusterContext()
	if err != nil {
		return errors.Wrap(errors.ErrorTypeCluster, "failed to load cluster context", err)
	}

	token, err := deps.GenerateK8sToken(ctx, "kubernetes-dashboard", "admin-user", 24)
	if err != nil {
		return errors.Wrap(errors.ErrorTypeK3s, "failed to generate bearer token", err)
	}

	localCmd := fmt.Sprintf("printf '%%s' \"%s\" | pbcopy", token)
	if _, err := deps.ExecuteLocalCommand(localCmd); err != nil {
		return errors.Wrap(errors.ErrorTypeSystem, "failed to copy bearer token to clipboard", err)
	}

	cmd.Println("Bearer token copied to clipboard.")
	return nil
}
