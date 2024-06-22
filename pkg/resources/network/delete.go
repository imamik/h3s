package network

import (
	"hcloud-k3s-cli/pkg/resources/clustercontext"
	"hcloud-k3s-cli/pkg/utils/logger"
)

func Delete(ctx clustercontext.ClusterContext) {
	network := Get(ctx)
	if network == nil {
		return
	}

	logger.LogResourceEvent(logger.Network, logger.Delete, network.Name, logger.Initialized)

	_, err := ctx.Client.Network.Delete(ctx.Context, network)
	if err != nil {
		logger.LogResourceEvent(logger.Network, logger.Delete, network.Name, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.Network, logger.Delete, network.Name, logger.Success)
}
