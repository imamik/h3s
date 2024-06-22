package cluster

import (
	"fmt"
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/cluster/network"
	"hcloud-k3s-cli/pkg/cluster/pool"
	"hcloud-k3s-cli/pkg/config"
)

func Create(conf config.Config) {
	fmt.Printf("Creating Cluster %s", conf.Name)

	ctx := clustercontext.Context(conf)

	ctx.Network = network.Create(ctx)
	pool.CreatePools(ctx)

	// Step 4: Create Firewall
	//err = createFirewall(ctx, client, conf, networkResp.IPRange)

	// Step 5: Create Load Balancer
	//err = createLoadBalancer(ctx, client, conf, networkResp)
}
