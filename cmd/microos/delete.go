package microos

import (
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/microos"
)

var Delete = &cobra.Command{
	Use:   "delete",
	Short: "Delete MicroOS snapshot/microos",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := clustercontext.Context()
		microos.Delete(ctx)
	},
}
