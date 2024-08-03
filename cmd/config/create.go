package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"h3s/internal/config/build"
	"h3s/internal/k3s/releases"
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
