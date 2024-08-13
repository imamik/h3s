package node

import (
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/utils/logger"
)

func Delete(
	ctx *cluster.Cluster,
	pool config.NodePool,
	i int,
) {
	server := Get(ctx, pool, i)

	if server == nil {
		return
	}

	addEvent, logEvents := logger.NewEventLogger(logger.Server, logger.Delete, server.Name)
	defer logEvents()

	_, _, err := ctx.CloudClient.Server.DeleteWithResult(ctx.Context, server)
	if err != nil {
		addEvent(logger.Failure, err)
		return
	}

	addEvent(logger.Success)
}
