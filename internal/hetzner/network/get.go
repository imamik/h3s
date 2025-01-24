package network

import (
	"errors"
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Get gets the Hetzner cloud network
func Get(ctx *cluster.Cluster) (*hcloud.Network, error) {
	networkName := getName(ctx)

	l := logger.New(nil, logger.Network, logger.Get, networkName)
	defer l.LogEvents()

	network, _, err := ctx.CloudClient.Network.GetByName(ctx.Context, networkName)
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
	return network, nil
}
