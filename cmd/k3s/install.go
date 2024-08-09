package k3s

import (
	"github.com/spf13/cobra"
	"h3s/internal/clustercontext"
	"h3s/internal/k3s/install"
)

var Install = &cobra.Command{
	Use:   "install",
	Short: "Install k3s on all servers in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := clustercontext.Context()
		install.Install(ctx)
		install.InstallSoftware(ctx)
		install.DownloadKubeconfig(ctx)
	},
}

var InstallSoftware = &cobra.Command{
	Use:   "software",
	Short: "Install software on all servers in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := clustercontext.Context()
		install.InstallSoftware(ctx)
	},
}

func init() {
	Install.AddCommand(InstallSoftware)
}
