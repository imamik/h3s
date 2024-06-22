package pool

import (
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/cluster/pool/placementgroup"
	"hcloud-k3s-cli/pkg/cluster/pool/server"
	"hcloud-k3s-cli/pkg/config"
)

func Delete(ctx clustercontext.ClusterContext) {
	// Delete control plane pool
	deletePool(ctx.Config.ControlPlane.Pool, ctx)

	// Delete worker pools
	for _, pool := range ctx.Config.WorkerPools {
		deletePool(pool, ctx)
	}
}

func deletePool(pool config.NodePool, ctx clustercontext.ClusterContext) {
	for i := 0; i < pool.Nodes; i++ {
		server.Delete(pool, i, ctx)
	}
	placementgroup.Delete(pool, ctx)
}
