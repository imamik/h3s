package kubectl

import (
	"h3s/cmd/dependencies"
	"h3s/cmd/errors"

	"github.com/spf13/cobra"
)

func runKubectl(cmd *cobra.Command, args []string) error {
	deps := dependencies.Get()

	ctx, err := deps.GetClusterContext()
	if err != nil {
		return errors.Wrap(errors.ErrorTypeCluster, "failed to load cluster context", err)
	}

	res, err := deps.ExecuteKubectlCommand(ctx, args)
	if err != nil {
		return errors.Wrap(errors.ErrorTypeKubectl, "failed to execute kubectl command", err)
	}

	cmd.Println(res)
	return nil
}
