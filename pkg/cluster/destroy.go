package cluster

import (
	"context"
	"hcloud-k3s-cli/pkg/cluster/network"
	"hcloud-k3s-cli/pkg/cluster/placementgroup"
	"hcloud-k3s-cli/pkg/cluster/utils"
	"hcloud-k3s-cli/pkg/config"
)

func Destroy(conf config.Config) {
	ctx := context.Background()
	client := utils.GetClient()

	_ = network.Delete(ctx, client, conf)
	_ = placementgroup.Delete(placementgroup.ControlPlanePool, ctx, client, conf)
}
