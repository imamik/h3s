package ssh

import (
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/ssh"
	"strings"
)

var Ssh = &cobra.Command{
	Use:                "ssh",
	Short:              "Proxy ssh commands to the first remote control plane server",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := clustercontext.Context()
		ssh.SSH(ctx, strings.Join(args, " "))
	},
}
