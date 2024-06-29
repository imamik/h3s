package pool

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/config"
	"hcloud-k3s-cli/internal/resources/microos"
	"hcloud-k3s-cli/internal/resources/pool/node"
	"hcloud-k3s-cli/internal/resources/pool/placementgroup"
	"hcloud-k3s-cli/internal/utils/logger"
	"sync"
)

func CreatePools(
	ctx clustercontext.ClusterContext,
	sshKey *hcloud.SSHKey,
	network *hcloud.Network,
) []*hcloud.Server {
	// Create control plane pool
	nodes := create(
		ctx,
		sshKey,
		network,
		ctx.Config.ControlPlane.Pool,
		true,
		ctx.Config.ControlPlane.AsWorkerPool,
	)

	// Create worker pools
	for _, pool := range ctx.Config.WorkerPools {
		workerNodes := create(
			ctx,
			sshKey,
			network,
			pool,
			false,
			true,
		)
		nodes = append(nodes, workerNodes...)
	}

	return nodes
}

func create(
	ctx clustercontext.ClusterContext,
	sshKey *hcloud.SSHKey,
	network *hcloud.Network,
	pool config.NodePool,
	isControlPlane bool,
	isWorker bool,
) []*hcloud.Server {
	logger.LogResourceEvent(logger.Pool, logger.Create, ctx.GetName(pool.Name), logger.Initialized)

	architecture := config.GetArchitecture(pool.Instance)
	// img := hcloud.Image{Name: "ubuntu-24.04"}
	img := microos.Create(ctx, architecture)

	placementGroup := placementgroup.Create(ctx, pool, isControlPlane, isWorker)

	// Create a channel to collect the nodes & setup a WaitGroup
	var nodes []*hcloud.Server
	nodeCh := make(chan *hcloud.Server, pool.Nodes)
	var wg sync.WaitGroup

	for i := 0; i < pool.Nodes; i++ {
		wg.Add(1) // Increment the WaitGroup counter
		go func(i int) {
			defer wg.Done() // Decrement the counter when the goroutine completes
			n := node.Create(
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
			nodeCh <- n // Send the created node to the channel
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
