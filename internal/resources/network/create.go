package network

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/ip"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Create(ctx clustercontext.ClusterContext) *hcloud.Network {
	network := Get(ctx)
	if network == nil {
		network = create(ctx)
	}
	return network
}

func create(ctx clustercontext.ClusterContext) *hcloud.Network {
	networkName := getName(ctx)

	addEvent, logEvents := logger.NewEventLogger(logger.Network, logger.Create, networkName)
	defer logEvents()

	network, _, err := ctx.Client.Network.Create(ctx.Context, hcloud.NetworkCreateOpts{
		Name:    networkName,
		IPRange: ip.GetIpRange("10.0.0.0/16"),
		Labels:  ctx.GetLabels(),
	})
	if err != nil || network == nil {
		addEvent(logger.Failure, err)
		return nil
	}

	addEvent(logger.Success)
	logEvents()
	addEvent, logEvents = logger.NewEventLogger(logger.Subnet, logger.Create, networkName)

	subnet := hcloud.NetworkSubnet{
		Type:        hcloud.NetworkSubnetTypeServer,
		IPRange:     ip.GetIpRange("10.0.0.0/16"),
		NetworkZone: ctx.Config.NetworkZone,
	}

	subNet, _, err := ctx.Client.Network.AddSubnet(ctx.Context, network, hcloud.NetworkAddSubnetOpts{
		Subnet: subnet,
	})
	if err != nil || subNet == nil {
		addEvent(logger.Failure, err)
		return nil
	}

	addEvent(logger.Success)
	return network
}
