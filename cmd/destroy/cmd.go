package destroy

import (
	"github.com/spf13/cobra"
)

// Cmd is the main command for destroying resources
var Cmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy various resources",
	Long:  `Destroy various resources - hetzner cloud resources, k8s components etc.`,
}

// destroyClusterCmd destroys an existing cluster (leaving the configuration files)
var destroyClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Destroy an existing cluster",
	Long:  `Destroy an existing cluster including all resources`,
	RunE:  runDestroyCluster,
}

// init adds subcommands to the main destroy command
func init() {
	Cmd.AddCommand(
		destroyClusterCmd,
	)
}
