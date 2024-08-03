package microos

import (
	"github.com/spf13/cobra"
	"h3s/internal/clustercontext"
	"h3s/internal/resources/microos"
)

var Delete = &cobra.Command{
	Use:   "delete",
	Short: "Delete MicroOS snapshot/microos",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := clustercontext.Context()
		microos.Delete(ctx)
	},
}
