package pool

import (
	"context"
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/pool/placementgroup"
	"hcloud-k3s-cli/pkg/cluster/pool/server"
	"hcloud-k3s-cli/pkg/config"
)

func Delete(conf config.Config, client *hcloud.Client, ctx context.Context) {
	// Delete control plane pool
	deletePool(conf.ControlPlane.Pool.Name, conf.ControlPlane.Pool.Nodes, conf, client, ctx)

	// Delete worker pools
	for _, pool := range conf.WorkerPools {
		deletePool(pool.Name, pool.Nodes, conf, client, ctx)
	}
}

func deletePool(
	name string,
	nodes int,
	conf config.Config,
	client *hcloud.Client,
	ctx context.Context,
) {

	for i := 0; i < nodes; i++ {
		nodeName := fmt.Sprintf("%s-node-%d", name, i+1)
		server.Delete(nodeName, conf, client, ctx)
	}

	placementgroup.Delete(name+"-pool", conf, client, ctx)
}
