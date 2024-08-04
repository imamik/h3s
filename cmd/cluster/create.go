package cluster

import (
	"github.com/spf13/cobra"
	"h3s/internal/clustercontext"
	"h3s/internal/k3s/install"
	"h3s/internal/resources/cluster"
)

// Create is a cobra.Command that handles the creation of a new cluster
var Create = &cobra.Command{
	Use:   "create",
	Short: "Create a new cluster",
	Long:  `Create a new cluster including the necessary resources, install k3s and configure the cluster`,
	Run:   runCreate,
}

func runCreate(cmd *cobra.Command, args []string) {
	// Get the cluster context
	ctx := clustercontext.Context()
	// Create the cluster resources
	cluster.Create(ctx)
	// Install k3s on the cluster
	install.Install(ctx)
}
