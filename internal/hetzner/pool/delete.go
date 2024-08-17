package pool

import (
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/hetzner/pool/node"
	"h3s/internal/hetzner/pool/placementgroup"
	"h3s/internal/utils/logger"
	"sync"
)

func Delete(ctx *cluster.Cluster) error {
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
	return nil
}

func deletePool(ctx *cluster.Cluster, pool config.NodePool) error {
	l := logger.New(nil, logger.Pool, logger.Delete, ctx.GetName(pool.Name))
	defer l.LogEvents()

	// Delete all nodes
	var wg sync.WaitGroup
	for i := 0; i < pool.Nodes; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			node.Delete(ctx, pool, i)
		}(i)
	}
	wg.Wait()

	// Delete placement group
	err := placementgroup.Delete(ctx, pool)
	if err != nil {
		return err
	}

	l.AddEvent(logger.Success)
	return nil
}
