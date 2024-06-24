package pool

import (
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/config"
	"hcloud-k3s-cli/internal/resources/pool/node"
	"hcloud-k3s-cli/internal/resources/pool/placementgroup"
	"hcloud-k3s-cli/internal/utils/logger"
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
		node.Delete(ctx, pool, i)
	}
	placementgroup.Delete(ctx, pool)
	logger.LogResourceEvent(logger.Pool, logger.Delete, ctx.GetName(pool.Name), logger.Success)
}
