package pool

import (
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/hetzner/pool/node"
	"h3s/internal/hetzner/pool/placementgroup"
	"h3s/internal/utils/logger"
	"sync"
)

// Delete removes the control plane and worker pools from the cluster.
func Delete(ctx *cluster.Cluster) error {
	l := logger.New(nil, logger.Pool, logger.Delete, ctx.GetName(ctx.Config.ControlPlane.Pool.Name))
	defer l.LogEvents()

	var wg sync.WaitGroup

	// Delete control plane pool
	wg.Add(1)
	go func() {
		if err := deletePool(ctx, ctx.Config.ControlPlane.Pool); err != nil {
			l.AddEvent(logger.Failure, err)
		}
		wg.Done()
	}()

	// Delete worker pools
	for _, pool := range ctx.Config.WorkerPools {
		wg.Add(1)
		go func(pool config.NodePool) {
			if err := deletePool(ctx, pool); err != nil {
				l.AddEvent(logger.Failure, err)
			}
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
			if err := node.Delete(ctx, pool, i); err != nil {
				l.AddEvent(logger.Failure, err)
			}
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
