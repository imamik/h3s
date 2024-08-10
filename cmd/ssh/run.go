package ssh

import (
	"fmt"
	"github.com/spf13/cobra"
	"h3s/internal/clustercontext"
	"h3s/internal/ssh"
	"strings"
)

// runSsh proxies ssh commands to the first control plane server in the h3s cluster
func runSsh(_ *cobra.Command, args []string) error {
	ctx := clustercontext.Context()

	cmd := strings.Join(args, " ")
	res, err := ssh.SSH(ctx, cmd)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}
