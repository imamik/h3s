package ssh

import (
	"fmt"
	"github.com/spf13/cobra"
	"h3s/internal/clustercontext"
	"h3s/internal/ssh"
	"strings"
)

// Ssh is the command to proxy ssh commands to the first remote control plane server
var Ssh = &cobra.Command{
	Use:                "ssh",
	Short:              "Proxy ssh commands to first control plane server",
	DisableFlagParsing: true,
	Run:                runSsh,
}

func runSsh(_ *cobra.Command, args []string) {
	ctx := clustercontext.Context()
	res, err := ssh.SSH(ctx, strings.Join(args, " "))
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
