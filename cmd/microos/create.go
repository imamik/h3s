package microos

import (
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/microos"
	"hcloud-k3s-cli/internal/resources/sshkey"
)

var Create = &cobra.Command{
	Use:   "create",
	Short: "Create MicroOS snapshot/microos",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := clustercontext.Context()
		sshkey.Create(ctx)
		microos.Create(ctx)
	},
}
