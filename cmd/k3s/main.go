package k3s

import (
	"github.com/spf13/cobra"
)

var K3s = &cobra.Command{
	Use:   "k3s",
	Short: "K3S utils",
}

func init() {
	K3s.AddCommand(Releases)
	K3s.AddCommand(Install)
}
