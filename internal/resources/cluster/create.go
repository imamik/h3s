package cluster

import (
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/dns"
	"hcloud-k3s-cli/internal/resources/gateway"
	"hcloud-k3s-cli/internal/resources/loadbalancers"
	"hcloud-k3s-cli/internal/resources/microos"
	"hcloud-k3s-cli/internal/resources/network"
	"hcloud-k3s-cli/internal/resources/pool"
	"hcloud-k3s-cli/internal/resources/server"
	"hcloud-k3s-cli/internal/resources/sshkey"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Create(ctx clustercontext.ClusterContext) {
	logger.LogResourceEvent(logger.Cluster, logger.Create, ctx.Config.Name, logger.Initialized)

	sshKey := sshkey.Create(ctx)
	net := network.Create(ctx)
	microos.Create(ctx)
	pool.CreatePools(ctx, sshKey, net)

	nodes := server.GetAll(ctx)
	loadbalancers.Create(ctx, net, nodes)

	lb := loadbalancers.Get(ctx)
	dns.Create(ctx, lb)

	if ctx.Config.PublicIps == false {
		gateway.Create(ctx)
	}

	logger.LogResourceEvent(logger.Cluster, logger.Create, ctx.Config.Name, logger.Success)
}
