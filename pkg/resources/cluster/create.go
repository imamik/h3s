package cluster

import (
	"hcloud-k3s-cli/pkg/clustercontext"
	"hcloud-k3s-cli/pkg/resources/network"
	"hcloud-k3s-cli/pkg/resources/pool"
	"hcloud-k3s-cli/pkg/resources/sshkey"
	"hcloud-k3s-cli/pkg/utils/logger"
)

func Create(ctx clustercontext.ClusterContext) {
	logger.LogResourceEvent(logger.Cluster, logger.Create, ctx.Config.Name, logger.Initialized)

	sshKey := sshkey.Create(ctx)
	net := network.Create(ctx)
	pool.CreatePools(ctx, sshKey, net)

	logger.LogResourceEvent(logger.Cluster, logger.Create, ctx.Config.Name, logger.Success)
}