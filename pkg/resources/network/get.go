package network

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/clustercontext"
	"hcloud-k3s-cli/pkg/utils/logger"
)

func Get(ctx clustercontext.ClusterContext) *hcloud.Network {
	networkName := getName(ctx)
	logger.LogResourceEvent(logger.Network, logger.Get, networkName, logger.Initialized)

	network, _, err := ctx.Client.Network.GetByName(ctx.Context, networkName)
	if err != nil || network == nil {
		logger.LogResourceEvent(logger.Network, logger.Get, networkName, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.Network, logger.Get, networkName, logger.Success)
	return network
}
