package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"h3s/internal/config/build"
	"h3s/internal/k3s/releases"
)

// Config is the command to configure project configuration for a h3s cluster
var Config = &cobra.Command{
	Use:   "create",
	Short: "Config project configuration",
	Long:  `Create a new h3s cluster configuration. This command will prompt you various questions to configure the project and create a new configuration file.`,
	Run:   runCreate,
}

func runCreate(cmd *cobra.Command, args []string) {
	k3sReleases, err := releases.GetFilteredReleases(false, true, 5)
	if err != nil {
		fmt.Println(err)
		return
	}

	build.Build(k3sReleases)
}
