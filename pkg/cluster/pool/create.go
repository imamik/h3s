package pool

import (
	"context"
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/pool/placementgroup"
	"hcloud-k3s-cli/pkg/cluster/pool/server"
	"hcloud-k3s-cli/pkg/config"
)

func CreatePools(n *hcloud.Network, conf config.Config, client *hcloud.Client, ctx context.Context) {
	// Create control plane pool
	create(conf.ControlPlane.Pool, true, conf.ControlPlane.AsWorkerPool, n, conf, client, ctx)

	// Create worker pools
	for _, pool := range conf.WorkerPools {
		create(pool, false, true, n, conf, client, ctx)
	}
}

func create(
	pool config.NodePool,
	isControlPlane bool,
	isWorker bool,
	network *hcloud.Network,
	conf config.Config,
	client *hcloud.Client,
	ctx context.Context,
) {
	placementGroup := placementgroup.Create(pool.Name+"-pool", conf, client, ctx)

	for i := 0; i < pool.Nodes; i++ {
		nodeName := fmt.Sprintf("%s-node-%d", pool.Name, i+1)
		server.Create(
			nodeName, isControlPlane, isWorker, pool.Location, pool.Instance,
			network, placementGroup, conf, client, ctx,
		)
	}
}
