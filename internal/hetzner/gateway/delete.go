package gateway

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/hetzner/network"
	"h3s/internal/utils/ip"
	"h3s/internal/utils/logger"
)

func Delete(ctx *cluster.Cluster) {
	net := network.Get(ctx)
	server, err := Get(ctx)
	if err != nil || server == nil {
		return
	}

	deleteGatewayRoute(ctx, net, server)
	delete(ctx, server)
}

func deleteGatewayRoute(ctx *cluster.Cluster, net *hcloud.Network, proxy *hcloud.Server) {
	if len(proxy.PrivateNet) == 0 {
		return
	}
	ctx.CloudClient.Network.DeleteRoute(ctx.Context, net, hcloud.NetworkDeleteRouteOpts{
		Route: hcloud.NetworkRoute{
			Destination: ip.GetIpRange("0.0.0.0/0"),
			Gateway:     proxy.PrivateNet[0].IP,
		},
	})
}

func delete(ctx *cluster.Cluster, server *hcloud.Server) {
	addEvent, logEvents := logger.NewEventLogger(logger.Server, logger.Delete, server.Name)
	defer logEvents()

	_, _, err := ctx.CloudClient.Server.DeleteWithResult(ctx.Context, server)
	if err != nil {
		addEvent(logger.Failure, err)
		return
	}

	addEvent(logger.Success)
}
