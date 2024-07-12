package kubectl

import (
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/internal/kubectl"
)

var Kubectl = &cobra.Command{
	Use:                "kubectl",
	Short:              "Proxy kubectl commands via ssh to the Kubernetes API server",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		kubectl.Execute(args)
	},
}
