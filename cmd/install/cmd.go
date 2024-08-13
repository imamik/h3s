package install

import (
	"github.com/spf13/cobra"
)

// Cmd is the main command for installing k3s or k8s components on a h3s cluster
var Cmd = &cobra.Command{
	Use:   "install",
	Short: "Install k3s or k8s components",
	Long:  `Install k3s or k8s components on a cluster`,
}

// installK3sCmd installs k3s on all servers in the h3s cluster
var installK3sCmd = &cobra.Command{
	Use:   "k3s",
	Short: "Install k3s on all servers in the cluster",
	RunE:  runInstallK3s,
}

// installComponentsCmd installs all components on the h3s cluster
var installComponentsCmd = &cobra.Command{
	Use:   "components",
	Short: "Install k8s components in the cluster",
	RunE:  runInstallComponents,
}

// init adds subcommands to the Cluster command
func init() {
	Cmd.AddCommand(
		installK3sCmd,
		installComponentsCmd,
	)
}
