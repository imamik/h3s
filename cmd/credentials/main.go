package credentials

import (
	"github.com/spf13/cobra"
	"h3s/internal/config/credentials"
)

// Credentials is the command to configure project credentials (Hetzner Cloud & DNS API token & k3s token)
var Credentials = &cobra.Command{
	Use:   "credentials",
	Short: "Configure project credentials",
	Run:   runCredentials,
}

func runCredentials(_ *cobra.Command, _ []string) {
	credentials.Configure()
}
