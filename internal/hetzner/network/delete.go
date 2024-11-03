package network

import (
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
)

// Delete removes the Hetzner cloud network
func Delete(ctx *cluster.Cluster) error {
	l := logger.New(nil, logger.Network, logger.Delete, getName(ctx))
	defer l.LogEvents()

	network, err := Get(ctx)
	if network == nil && err.Error() == "network is nil" {
		l.AddEvent(logger.Success, "no network found to delete")
		return nil
	}

	_, err = ctx.CloudClient.Network.Delete(ctx.Context, network)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	l.AddEvent(logger.Success)
	return nil
}
