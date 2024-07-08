package cluster

import (
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/loadbalancers"
	"hcloud-k3s-cli/internal/resources/network"
	"hcloud-k3s-cli/internal/resources/pool"
	"hcloud-k3s-cli/internal/resources/proxy"
	"hcloud-k3s-cli/internal/resources/sshkey"
	"hcloud-k3s-cli/internal/utils/file"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Destroy(ctx clustercontext.ClusterContext) {
	logger.LogResourceEvent(logger.Cluster, logger.Delete, ctx.Config.Name, logger.Initialized)

	loadbalancers.Delete(ctx)
	pool.Delete(ctx)
	network.Delete(ctx)
	sshkey.Delete(ctx)
	proxy.Delete(ctx)

	file.Delete("./k3s.yaml")

	logger.LogResourceEvent(logger.Cluster, logger.Delete, ctx.Config.Name, logger.Success)
}
