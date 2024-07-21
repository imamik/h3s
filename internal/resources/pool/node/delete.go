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

	addEvent, logEvents := logger.NewEventLogger(logger.Server, logger.Delete, server.Name)
	defer logEvents()

	_, _, err := ctx.Client.Server.DeleteWithResult(ctx.Context, server)
	if err != nil {
		addEvent(logger.Failure, err)
		return
	}

	addEvent(logger.Success)
}
