package network

import (
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Delete(ctx clustercontext.ClusterContext) {
	network := Get(ctx)
	if network == nil {
		return
	}

	addEvent, logEvents := logger.NewEventLogger(logger.Network, logger.Delete, network.Name)
	defer logEvents()

	_, err := ctx.Client.Network.Delete(ctx.Context, network)
	if err != nil {
		addEvent(logger.Failure, err)
		return
	}

	addEvent(logger.Success)
}
