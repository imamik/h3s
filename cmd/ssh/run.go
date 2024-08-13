package ssh

import (
	"github.com/spf13/cobra"
	"h3s/internal/cluster"
	"h3s/internal/utils/common"
	"strings"
)

// runSsh proxies ssh commands to the first control plane server in the h3s cluster
func runSsh(cmd *cobra.Command, args []string) error {
	ctx, err := cluster.Context()

	if err != nil {
		return err
	}

	command := strings.Join(args, " ")
	res, err := common.SSH(ctx, command)
	if err != nil {
		return err
	}

	cmd.Println(res)
	return nil
}
