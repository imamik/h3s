package pool

import (
	"hcloud-k3s-cli/pkg/clustercontext"
	"hcloud-k3s-cli/pkg/config"
	"hcloud-k3s-cli/pkg/resources/placementgroup"
	"hcloud-k3s-cli/pkg/resources/server"
	"hcloud-k3s-cli/pkg/utils/logger"
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
	logger.LogResourceEvent(logger.Pool, logger.Delete, ctx.GetName(pool.Name), logger.Initialized)
	for i := 0; i < pool.Nodes; i++ {
		server.Delete(ctx, pool, i)
	}
	placementgroup.Delete(ctx, pool)
	logger.LogResourceEvent(logger.Pool, logger.Delete, ctx.GetName(pool.Name), logger.Success)
}
