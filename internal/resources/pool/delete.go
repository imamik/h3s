package pool

import (
	"h3s/internal/clustercontext"
	"h3s/internal/config"
	"h3s/internal/resources/pool/node"
	"h3s/internal/resources/pool/placementgroup"
	"h3s/internal/utils/logger"
	"sync"
)

func Delete(ctx clustercontext.ClusterContext) {
	var wg sync.WaitGroup

	// Delete control plane pool
	wg.Add(1)
	go func() {
		deletePool(ctx, ctx.Config.ControlPlane.Pool)
		wg.Done()
	}()

	// Delete worker pools
	for _, pool := range ctx.Config.WorkerPools {
		wg.Add(1)
		go func(pool config.NodePool) {
			deletePool(ctx, pool)
			wg.Done()
		}(pool)
	}

	wg.Wait()
}

func deletePool(ctx clustercontext.ClusterContext, pool config.NodePool) {
	addEvent, logEvents := logger.NewEventLogger(logger.Pool, logger.Delete, ctx.GetName(pool.Name))
	defer logEvents()

	var wg sync.WaitGroup
	for i := 0; i < pool.Nodes; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			node.Delete(ctx, pool, i)
		}(i)
	}
	wg.Wait()
	placementgroup.Delete(ctx, pool)
	addEvent(logger.Success)
}
