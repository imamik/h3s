package cluster

import (
	"context"
	"hcloud-k3s-cli/pkg/cluster/network"
	"hcloud-k3s-cli/pkg/cluster/pool"
	"hcloud-k3s-cli/pkg/cluster/utils"
	"hcloud-k3s-cli/pkg/config"
)

func Create(conf config.Config) {
	ctx := context.Background()
	client := utils.GetClient()

	n := network.Create(conf, client, ctx)

	pool.CreatePools(n, conf, client, ctx)

	// Step 4: Create Firewall
	//err = createFirewall(ctx, client, conf, networkResp.IPRange)

	// Step 5: Create Load Balancer
	//err = createLoadBalancer(ctx, client, conf, networkResp)
}
