package network

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/clustercontext"
	"hcloud-k3s-cli/pkg/utils/logger"
	"net"
)

func Create(ctx clustercontext.ClusterContext) *hcloud.Network {
	network := Get(ctx)
	if network == nil {
		network = create(ctx)
	}
	return network
}

func getIpRange(s string) *net.IPNet {
	_, ipRange, err := net.ParseCIDR(s)
	if err != nil {
		logger.LogError("Invalid IP Range", err)
	}
	return ipRange
}

func create(ctx clustercontext.ClusterContext) *hcloud.Network {
	networkName := getName(ctx)

	logger.LogResourceEvent(logger.Network, logger.Create, networkName, logger.Initialized)

	network, _, err := ctx.Client.Network.Create(ctx.Context, hcloud.NetworkCreateOpts{
		Name:    networkName,
		IPRange: getIpRange("10.0.0.0/16"),
		Labels:  ctx.GetLabels(),
	})
	if err != nil || network == nil {
		logger.LogResourceEvent(logger.Network, logger.Create, networkName, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.Network, logger.Create, networkName, logger.Success)
	logger.LogResourceEvent(logger.Subnet, logger.Create, networkName, logger.Initialized)

	subnet := hcloud.NetworkSubnet{
		Type:        hcloud.NetworkSubnetTypeServer,
		IPRange:     getIpRange("10.0.0.0/16"),
		NetworkZone: ctx.Config.NetworkZone,
	}

	subNet, _, err := ctx.Client.Network.AddSubnet(ctx.Context, network, hcloud.NetworkAddSubnetOpts{
		Subnet: subnet,
	})
	if err != nil || subNet == nil {
		logger.LogResourceEvent(logger.Network, logger.Create, networkName, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.Network, logger.Create, networkName, logger.Success)

	return network
}
