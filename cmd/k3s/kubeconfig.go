package k3s

import (
	"github.com/spf13/cobra"
	"h3s/internal/clustercontext"
	"h3s/internal/k3s/install"
	"h3s/internal/k3s/kubeconfig"
)

var KubeConfig = &cobra.Command{
	Use:   "kubeconfig",
	Short: "Get the kubeconfig for the k3s cluster",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := clustercontext.Context()
		_, _, gatewayServer, controlPlaneNodes, _ := install.GetSetup(ctx)
		kubeconfig.Download(ctx, gatewayServer, controlPlaneNodes[0])
	},
}