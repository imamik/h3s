package kubectl

import (
	"github.com/spf13/cobra"
)

// Cmd proxies kubectl commands either directly with the kubeconfig if available or via SSH to the first control plane server
var Cmd = &cobra.Command{
	Use:                "kubectl",
	Short:              "Run kubectl commands",
	Long:               `Run kubectl commands either directly (if setup and possible) or via SSH to the first control plane server`,
	DisableFlagParsing: true,
	RunE:               runKubectl,
}
