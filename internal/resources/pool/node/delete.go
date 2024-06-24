package node

import (
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/config"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Delete(
	ctx clustercontext.ClusterContext,
	pool config.NodePool,
	i int,
) {
	server := Get(ctx, pool, i)

	if server == nil {
		return
	}

	logger.LogResourceEvent(logger.Server, logger.Delete, server.Name, logger.Initialized)

	_, _, err := ctx.Client.Server.DeleteWithResult(ctx.Context, server)
	if err != nil {
		logger.LogResourceEvent(logger.Server, logger.Delete, server.Name, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.Server, logger.Delete, server.Name, logger.Success)
}
