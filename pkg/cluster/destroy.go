package cluster

import (
	"context"
	"hcloud-k3s-cli/pkg/cluster/network"
	"hcloud-k3s-cli/pkg/cluster/pool"
	"hcloud-k3s-cli/pkg/cluster/utils"
	"hcloud-k3s-cli/pkg/config"
)

func Destroy(conf config.Config) {
	ctx := context.Background()
	client := utils.GetClient()

	network.Delete(ctx, client, conf)
	pool.Delete(conf, client, ctx)
}
