package microos

import (
	"github.com/spf13/cobra"
)

var Image = &cobra.Command{
	Use:   "microos",
	Short: "Utils for MicroOS snapshots/images",
}

func init() {
	Image.AddCommand(Create)
	Image.AddCommand(Delete)
}
