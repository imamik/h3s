package pool

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/config"
	"hcloud-k3s-cli/internal/resources/microos/image"
	"hcloud-k3s-cli/internal/resources/network"
	"hcloud-k3s-cli/internal/resources/pool/node"
	"hcloud-k3s-cli/internal/resources/pool/placementgroup"
	"hcloud-k3s-cli/internal/resources/sshkey"
	"hcloud-k3s-cli/internal/utils/logger"
	"sync"
)

func CreatePools(ctx clustercontext.ClusterContext) []*hcloud.Server {
	sshKey := sshkey.Get(ctx)
	net := network.Get(ctx)

	// Create a channel to collect the nodes & setup a WaitGroup
	nodeCh := make(chan []*hcloud.Server)
	var wg sync.WaitGroup

	// Create control plane pool in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		nodeCh <- CreatePool(
			ctx,
			sshKey,
			net,
			ctx.Config.ControlPlane.Pool,
			true,
			ctx.Config.ControlPlane.AsWorkerPool,
		)
	}()

	// Create worker pools in separate goroutines
	for _, pool := range ctx.Config.WorkerPools {
		wg.Add(1)
		go func(pool config.NodePool) {
			defer wg.Done()
			nodeCh <- CreatePool(
				ctx,
				sshKey,
				net,
				pool,
				false,
				true,
			)
		}(pool)
	}

	// Wait for all goroutines to finish & close the channel
	go func() {
		wg.Wait()
		close(nodeCh)
	}()

	// Collect the nodes from the channel
	var nodes []*hcloud.Server
	for n := range nodeCh {
		nodes = append(nodes, n...)
	}

	return nodes
}

func CreatePool(
	ctx clustercontext.ClusterContext,
	sshKey *hcloud.SSHKey,
	network *hcloud.Network,
	pool config.NodePool,
	isControlPlane bool,
	isWorker bool,
) []*hcloud.Server {
	logger.LogResourceEvent(logger.Pool, logger.Create, ctx.GetName(pool.Name), logger.Initialized)

	img, err := image.Get(ctx, config.GetArchitecture(pool.Instance))
	if err != nil {
		logger.LogResourceEvent(logger.Pool, logger.Create, ctx.GetName(pool.Name), logger.Failure)
		return nil
	}

	placementGroup := placementgroup.Create(ctx, pool, isControlPlane, isWorker)

	// Create a channel to collect the nodes & setup a WaitGroup
	var nodes []*hcloud.Server
	nodeCh := make(chan *hcloud.Server, pool.Nodes)
	var wg sync.WaitGroup

	for i := 0; i < pool.Nodes; i++ {
		wg.Add(1) // Increment the WaitGroup counter
		go func(i int) {
			defer wg.Done() // Decrement the counter when the goroutine completes
			nodeCh <- node.Create(
				ctx,
				sshKey,
				network,
				img,
				placementGroup,
				pool,
				i,
				isControlPlane,
				isWorker,
			)
		}(i)
	}

	// Wait for all goroutines to finish & close the channel
	wg.Wait()
	close(nodeCh)

	// Collect the nodes from the channel
	for n := range nodeCh {
		nodes = append(nodes, n)
	}

	logger.LogResourceEvent(logger.Pool, logger.Create, ctx.GetName(pool.Name), logger.Success)
	return nodes
}
