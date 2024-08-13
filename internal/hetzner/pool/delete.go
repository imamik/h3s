package pool

import (
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/hetzner/pool/node"
	"h3s/internal/hetzner/pool/placementgroup"
	"h3s/internal/utils/logger"
	"sync"
)

func Delete(ctx *cluster.Cluster) {
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

func deletePool(ctx *cluster.Cluster, pool config.NodePool) {
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
