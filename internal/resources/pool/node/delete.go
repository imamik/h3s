package node

import (
	"h3s/internal/clustercontext"
	"h3s/internal/config"
	"h3s/internal/utils/logger"
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
