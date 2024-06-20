package cfg

import (
	"fmt"
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/pkg/config/build"
	"hcloud-k3s-cli/pkg/k3s/releases"
)

var Create = &cobra.Command{
	Use:   "create",
	Short: "Create project configuration",
	Run: func(cmd *cobra.Command, args []string) {

		k3sReleases, err := releases.GetFilteredReleases(false, true, 5)
		if err != nil {
			fmt.Println(err)
			return
		}

		build.Build(k3sReleases)
	},
}
