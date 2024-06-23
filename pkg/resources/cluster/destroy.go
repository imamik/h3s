package cluster

import (
	"hcloud-k3s-cli/pkg/clustercontext"
	"hcloud-k3s-cli/pkg/resources/network"
	"hcloud-k3s-cli/pkg/resources/pool"
	"hcloud-k3s-cli/pkg/resources/sshkey"
	"hcloud-k3s-cli/pkg/utils/logger"
)

func Destroy(ctx clustercontext.ClusterContext) {
	logger.LogResourceEvent(logger.Cluster, logger.Delete, ctx.Config.Name, logger.Initialized)

	pool.Delete(ctx)
	network.Delete(ctx)
	sshkey.Delete(ctx)

	logger.LogResourceEvent(logger.Cluster, logger.Delete, ctx.Config.Name, logger.Success)
}
