package image

import (
	"github.com/spf13/cobra"
)

var (
	all bool
	arm bool
	x86 bool
	l   string
)

var Image = &cobra.Command{
	Use:   "image",
	Short: "Utils for MicroOS snapshots/images",
}

func init() {
	Image.AddCommand(Create)
	Image.AddCommand(Delete)
}
