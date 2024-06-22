package config

import (
	"github.com/spf13/cobra"
)

var Config = &cobra.Command{
	Use:   "config",
	Short: "Utils for configuration",
}

func init() {
	Config.AddCommand(Create)
	Config.AddCommand(Credentials)
}
