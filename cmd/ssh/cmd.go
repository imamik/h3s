package ssh

import (
	"github.com/spf13/cobra"
)

// Cmd is the command to proxy ssh commands to the first remote control plane server
var Cmd = &cobra.Command{
	Use:                "ssh",
	Short:              "Proxy ssh commands to first control plane server",
	DisableFlagParsing: true,
	RunE:               runSsh,
}
