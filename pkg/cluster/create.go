package cluster

import (
	"context"
	"hcloud-k3s-cli/pkg/cluster/network"
	"hcloud-k3s-cli/pkg/cluster/placementgroup"
	"hcloud-k3s-cli/pkg/cluster/utils"
	"hcloud-k3s-cli/pkg/config"
)

func Create(conf config.Config) {
	ctx := context.Background()
	client := utils.GetClient()

	network.Create(ctx, client, conf)

	placementgroup.Create(placementgroup.ControlPlanePool, ctx, client, conf)

	// Step 3: Create Servers
	//err = createServers(ctx, client, conf, networkResp, placementGroupResp)

	// Step 4: Create Firewall
	//err = createFirewall(ctx, client, conf, networkResp.IPRange)

	// Step 5: Create Load Balancer
	//err = createLoadBalancer(ctx, client, conf, networkResp)
}
