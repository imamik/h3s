package gateway

import (
	"errors"
	"h3s/internal/cluster"
	"h3s/internal/hetzner/network"
	"h3s/internal/utils/ip"
	"h3s/internal/utils/logger"
	"net"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Delete deletes the Hetzner cloud gateway
func Delete(ctx *cluster.Cluster) error {
	// Get network
	network, err := network.Get(ctx)
	if err != nil {
		return err
	}

	// Get gateway server
	server, err := Get(ctx)
	if err != nil {
		return err
	}
	if server == nil {
		return errors.New("no gateway server found")
	}

	// Delete gateway route
	if err := deleteGatewayRoute(ctx, network, server); err != nil {
		return err
	}

	return deleteGateway(ctx, server)
}

func deleteGatewayRoute(ctx *cluster.Cluster, network *hcloud.Network, proxy *hcloud.Server) error {
	if len(proxy.PrivateNet) == 0 {
		return errors.New("no private network found")
	}
	_, ipRange, err := net.ParseCIDR("0.0.0.0/0")
	if err != nil {
		return err
	}

	_, _, err = ctx.CloudClient.Network.DeleteRoute(ctx.Context, network, hcloud.NetworkDeleteRouteOpts{
		Route: hcloud.NetworkRoute{
			Destination: ipRange,
			Gateway:     ip.Private(proxy),
		},
	})
	return err
}

func deleteGateway(ctx *cluster.Cluster, server *hcloud.Server) error {
	l := logger.New(nil, logger.Server, logger.Delete, server.Name)
	defer l.LogEvents()

	// Delete server
	_, _, err := ctx.CloudClient.Server.DeleteWithResult(ctx.Context, server)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	l.AddEvent(logger.Success)
	return nil
}
