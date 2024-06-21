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
	create("control-plane", true, conf.ControlPlane.AsWorkerPool, conf.ControlPlane.Location, conf.ControlPlane.Nodes, n, conf, client, ctx)

	// Create worker pools
	for _, pool := range conf.WorkerPools {
		create(pool.Name, false, true, pool.Location, pool.Nodes, n, conf, client, ctx)
	}
}

func create(
	name string,
	isControlPlane bool,
	isWorker bool,
	location config.Location,
	nodes int,
	network *hcloud.Network,
	conf config.Config,
	client *hcloud.Client,
	ctx context.Context,
) {
	placementGroup := placementgroup.Create(name+"-pool", conf, client, ctx)

	for i := 0; i < nodes; i++ {
		nodeName := fmt.Sprintf("%s-node-%d", name, i+1)
		server.Create(nodeName, isControlPlane, isWorker, location, network, placementGroup, conf, client, ctx)
	}
}
