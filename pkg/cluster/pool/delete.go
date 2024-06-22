package pool

import (
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/cluster/pool/placementgroup"
	"hcloud-k3s-cli/pkg/cluster/pool/server"
	"hcloud-k3s-cli/pkg/config"
	"log"
)

func Delete(ctx clustercontext.ClusterContext) {
	// Delete control plane pool
	deletePool(ctx, ctx.Config.ControlPlane.Pool)

	// Delete worker pools
	for _, pool := range ctx.Config.WorkerPools {
		deletePool(ctx, pool)
	}
}

func deletePool(ctx clustercontext.ClusterContext, pool config.NodePool) {
	log.Println("Deleting pool", pool.Name)
	for i := 0; i < pool.Nodes; i++ {
		server.Delete(ctx, pool, i)
	}
	placementgroup.Delete(ctx, pool)
	log.Println("Deleted pool", pool.Name)
}
