// Package ssh provides functionality to proxy ssh commands to the first control plane server in the h3s cluster.
package ssh

import (
	"github.com/spf13/cobra"
)

// Cmd proxies ssh commands to the first remote control plane server in the h3s cluster
var Cmd = &cobra.Command{
	Use:                "ssh",
	Short:              "Proxy ssh commands to first control plane server",
	DisableFlagParsing: true,
	RunE:               runSSH,
}
