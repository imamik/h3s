package image

import (
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/config"
	"hcloud-k3s-cli/internal/resources/image"
	"log"
	"sync"
)

var Create = &cobra.Command{
	Use:   "create",
	Short: "Create MicroOS snapshot/image",
	Run: func(cmd *cobra.Command, args []string) {
		if !arm && !x86 && !all {
			log.Fatalf("Please specify at least one architecture to delete by using --arm or --all or both")
		}

		ctx := clustercontext.Context()

		var wg sync.WaitGroup

		loc := config.Location(l)

		if arm || all {
			wg.Add(1)
			go func() {
				defer wg.Done()
				image.Create(ctx, "arm", loc)
			}()
		}

		if x86 || all {
			wg.Add(1)
			go func() {
				defer wg.Done()
				image.Create(ctx, "x86", loc)
			}()
		}

		wg.Wait()
	},
}

func init() {
	Create.Flags().BoolVar(&arm, "arm", false, "")
	Create.Flags().BoolVar(&x86, "x86", false, "")
	Create.Flags().StringVar(&l, "location", string(config.Nuernberg), "")
}
