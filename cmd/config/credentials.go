package config

import (
	"github.com/spf13/cobra"
	"h3s/internal/config/credentials"
)

var Credentials = &cobra.Command{
	Use:   "credentials",
	Short: "Configure project credentials",
	Run: func(cmd *cobra.Command, args []string) {
		credentials.Configure()
	},
}
