package network

import (
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
)

func Delete(ctx *cluster.Cluster) {
	network := Get(ctx)
	if network == nil {
		return
	}

	addEvent, logEvents := logger.NewEventLogger(logger.Network, logger.Delete, network.Name)
	defer logEvents()

	_, err := ctx.CloudClient.Network.Delete(ctx.Context, network)
	if err != nil {
		addEvent(logger.Failure, err)
		return
	}

	addEvent(logger.Success)
}
