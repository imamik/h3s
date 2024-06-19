package config

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/pkg/config/build"
)

var InitCmd = &cobra.Command{
	Use:   "config build",
	Short: "Initialize project configuration",
	Run: func(cmd *cobra.Command, args []string) {

		config, err := build.InitConfig()
		if err != nil {
			fmt.Println(err)
			return
		}

		err = build.Save(config, "config.yaml")
		if err != nil {
			fmt.Println(err)
			return
		}

		color.Green("🎉🎉🎉 Project configuration 🛠️ initialized successfully 🎉🎉🎉")
	},
}
