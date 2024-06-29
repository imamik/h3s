package microos

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/microos"
	"log"
	"sync"
)

var Create = &cobra.Command{
	Use:   "create",
	Short: "Create MicroOS snapshot/microos",
	Run: func(cmd *cobra.Command, args []string) {
		if !arm && !x86 && !all {
			log.Fatalf("Please specify at least one architecture to delete by using --arm or --all or both")
		}

		ctx := clustercontext.Context()

		var wg sync.WaitGroup

		if arm || all {
			wg.Add(1)
			go func() {
				defer wg.Done()
				microos.Create(ctx, hcloud.ArchitectureARM)
			}()
		}

		if x86 || all {
			wg.Add(1)
			go func() {
				defer wg.Done()
				microos.Create(ctx, hcloud.ArchitectureX86)
			}()
		}

		wg.Wait()
	},
}

func init() {
	Create.Flags().BoolVar(&arm, "arm", false, "")
	Create.Flags().BoolVar(&x86, "x86", false, "")
	Create.Flags().BoolVar(&all, "all", false, "")
}
