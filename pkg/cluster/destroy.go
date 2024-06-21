package cluster

import (
	"context"
	"hcloud-k3s-cli/pkg/cluster/network"
	placementgroup2 "hcloud-k3s-cli/pkg/cluster/pool/placementgroup"
	"hcloud-k3s-cli/pkg/cluster/utils"
	"hcloud-k3s-cli/pkg/config"
)

func Destroy(conf config.Config) {
	ctx := context.Background()
	client := utils.GetClient()

	network.Delete(ctx, client, conf)
	placementgroup2.Delete(placementgroup2.ControlPlanePool, ctx, client, conf)
}
