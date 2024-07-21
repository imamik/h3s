package network

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Get(ctx clustercontext.ClusterContext) *hcloud.Network {
	networkName := getName(ctx)

	addEvent, logEvents := logger.NewEventLogger(logger.Network, logger.Get, networkName)
	defer logEvents()

	network, _, err := ctx.Client.Network.GetByName(ctx.Context, networkName)
	if err != nil || network == nil {
		addEvent(logger.Failure, err)
		return nil
	}

	addEvent(logger.Success)
	return network
}
