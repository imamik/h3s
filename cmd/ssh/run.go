package ssh

import (
	"fmt"
	"github.com/spf13/cobra"
	"h3s/internal/clustercontext"
	"h3s/internal/ssh"
	"strings"
)

// runSsh is the function that is executed when the ssh command is called - it proxies ssh commands to the first control plane server
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
