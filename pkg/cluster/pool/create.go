package pool

import (
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/cluster/pool/placementgroup"
	"hcloud-k3s-cli/pkg/cluster/pool/server"
	"hcloud-k3s-cli/pkg/config"
)

func CreatePools(ctx clustercontext.ClusterContext) {
	// Create control plane pool
	create(ctx.Config.ControlPlane.Pool, true, ctx.Config.ControlPlane.AsWorkerPool, ctx)

	// Create worker pools
	for _, pool := range ctx.Config.WorkerPools {
		create(pool, false, true, ctx)
	}
}

func create(
	pool config.NodePool,
	isControlPlane bool,
	isWorker bool,
	ctx clustercontext.ClusterContext,
) {
	placementGroup := placementgroup.Create(pool, ctx)

	for i := 0; i < pool.Nodes; i++ {
		server.Create(pool, i, isControlPlane, isWorker, placementGroup, ctx)
	}
}
