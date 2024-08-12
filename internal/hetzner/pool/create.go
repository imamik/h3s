package pool

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/hetzner/microos/image"
	"h3s/internal/hetzner/network"
	"h3s/internal/hetzner/pool/node"
	"h3s/internal/hetzner/pool/placementgroup"
	"h3s/internal/hetzner/sshkey"
	"h3s/internal/utils/logger"
	"sync"
)

func CreatePools(ctx *cluster.Cluster) []*hcloud.Server {
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
	ctx *cluster.Cluster,
	sshKey *hcloud.SSHKey,
	network *hcloud.Network,
	pool config.NodePool,
	isControlPlane bool,
	isWorker bool,
) []*hcloud.Server {
	addEvent, logEvents := logger.NewEventLogger(logger.Pool, logger.Create, ctx.GetName(pool.Name))
	defer logEvents()

	var img *hcloud.Image
	var err error
	if ctx.Config.Image == config.ImageMicroOS {
		img, err = image.Get(ctx, config.GetArchitecture(pool.Instance))
		if err != nil {
			addEvent(logger.Failure, err)
			return nil
		}
	} else {
		img = &hcloud.Image{
			Name: string(ctx.Config.Image),
		}
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

	addEvent(logger.Success)
	return nodes
}
