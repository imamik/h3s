// Package network contains the functionality for creating a Hetzner cloud network
package network

import (
	"errors"
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
	"net"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Create creates the Hetzner cloud network
func Create(ctx *cluster.Cluster) (*hcloud.Network, error) {
	network, err := Get(ctx)
	if network != nil && err == nil {
		return network, nil
	}
	return create(ctx)
}

func create(ctx *cluster.Cluster) (*hcloud.Network, error) {
	networkName := getName(ctx)

	l := logger.New(nil, logger.Network, logger.Create, networkName)
	defer l.LogEvents()

	_, ipRange, err := net.ParseCIDR("10.0.0.0/16")
	if err != nil {
		return nil, err
	}
	network, _, err := ctx.CloudClient.Network.Create(ctx.Context, hcloud.NetworkCreateOpts{
		Name:    networkName,
		IPRange: ipRange,
		Labels:  ctx.GetLabels(),
	})
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}
	if network == nil {
		err = errors.New("network is nil")
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	l.AddEvent(logger.Success)
	l.LogEvents()
	l = logger.New(nil, logger.Subnet, logger.Create, networkName)

	subnet := hcloud.NetworkSubnet{
		Type:        hcloud.NetworkSubnetTypeCloud,
		IPRange:     ipRange,
		NetworkZone: ctx.Config.NetworkZone,
	}

	subNet, _, err := ctx.CloudClient.Network.AddSubnet(ctx.Context, network, hcloud.NetworkAddSubnetOpts{
		Subnet: subnet,
	})
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}
	if subNet == nil {
		err = errors.New("subnet is nil")
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	l.AddEvent(logger.Success)
	return network, nil
}
