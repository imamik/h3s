package ssh

import (
	"github.com/spf13/cobra"
	"h3s/internal/clustercontext"
	"h3s/internal/utils/ssh"
	"strings"
)

// runSsh proxies ssh commands to the first control plane server in the h3s cluster
func runSsh(cmd *cobra.Command, args []string) error {
	ctx := clustercontext.Context()

	command := strings.Join(args, " ")
	res, err := ssh.SSH(ctx, command)
	if err != nil {
		return err
	}

	cmd.Println(res)
	return nil
}
