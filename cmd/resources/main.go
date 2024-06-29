package resources

import (
	"github.com/spf13/cobra"
)

var Resources = &cobra.Command{
	Use:   "resources",
	Short: "CLI to manage resources on Hetzner Cloud",
}

func init() {
	Resources.AddCommand(Create)
	Resources.AddCommand(Destroy)
}
