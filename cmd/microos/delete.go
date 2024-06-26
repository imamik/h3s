package microos

import (
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/microos"
	"log"
)

var Delete = &cobra.Command{
	Use:   "delete",
	Short: "Delete MicroOS snapshot/microos",
	Run: func(cmd *cobra.Command, args []string) {
		if !arm && !x86 && !all {
			log.Fatalf("Please specify at least one architecture to delete by using --arm or --all or both")
		}

		ctx := clustercontext.Context()

		if arm || all {
			microos.Delete(ctx, "arm")
		}
		if x86 || all {
			microos.Delete(ctx, "x86")
		}
	},
}

func init() {
	Delete.Flags().BoolVar(&arm, "arm", false, "")
	Delete.Flags().BoolVar(&x86, "x86", false, "")
	Delete.Flags().BoolVar(&all, "all", false, "")
}
