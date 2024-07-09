package k3s

import (
	"github.com/spf13/cobra"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/k3s/install"
)

var cleanup bool

var Install = &cobra.Command{
	Use:   "install",
	Short: "Install k3s on all servers in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := clustercontext.Context()
		install.Install(ctx, cleanup)
	},
}

var InstallSoftware = &cobra.Command{
	Use:   "software",
	Short: "Install software on all servers in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := clustercontext.Context()
		install.InstallSoftware(ctx, cleanup)
	},
}

func init() {
	InstallSoftware.Flags().BoolVar(&cleanup, "cleanup", true, "Force installation")
	Install.Flags().BoolVar(&cleanup, "cleanup", true, "Force installation")
	Install.AddCommand(InstallSoftware)
}
