package cfg

import (
	"github.com/spf13/cobra"
)

var Cfg = &cobra.Command{
	Use:   "cfg",
	Short: "Utils for configuration",
}

func init() {
	Cfg.AddCommand(Create)
}
