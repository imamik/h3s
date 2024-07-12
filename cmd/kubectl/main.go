package kubectl

import (
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/kubectl"
)

var Kubectl = &cobra.Command{
	Use:                "kubectl",
	Short:              "Proxy kubectl commands via ssh to the Kubernetes API server",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		useSsh := false
		var filteredArgs []string

		for _, arg := range args {
			if arg == "--ssh" {
				useSsh = true
			} else {
				filteredArgs = append(filteredArgs, arg)
			}
		}

		if useSsh {
			ctx := clustercontext.Context()
			err := kubectl.SSH(ctx, filteredArgs)
			if err != nil {
				panic(err)
			}
		} else {
			kubectl.Kubectl(filteredArgs)
		}
	},
}
