package ssh

import (
	"h3s/cmd/dependencies"
	"h3s/cmd/errors"
	"strings"

	"github.com/spf13/cobra"
)

// runSsh proxies ssh commands to the first control plane server in the h3s cluster
func runSsh(cmd *cobra.Command, args []string) error {
	deps := dependencies.Get()

	ctx, err := deps.GetClusterContext()
	if err != nil {
		return errors.Wrap(errors.ErrorTypeCluster, "failed to load cluster context", err)
	}

	command := strings.Join(args, " ")
	res, err := deps.ExecuteSSHCommand(ctx, command)
	if err != nil {
		return errors.Wrap(errors.ErrorTypeSSH, "failed to execute ssh command", err)
	}

	cmd.Println(res)
	return nil
}
