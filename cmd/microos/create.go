package microos

import (
	"github.com/spf13/cobra"
	"h3s/internal/clustercontext"
	"h3s/internal/resources/microos"
	"h3s/internal/resources/sshkey"
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
