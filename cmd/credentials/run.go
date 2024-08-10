package credentials

import (
	"github.com/spf13/cobra"
	"h3s/internal/config/credentials"
)

// runCreate is the function that is executed when the create command is called - it creates a new h3s cluster configuration
func runCreate(_ *cobra.Command, _ []string) error {
	credentials.Configure()
	return nil
}
