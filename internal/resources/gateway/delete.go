package gateway

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/clustercontext"
	"h3s/internal/resources/network"
	"h3s/internal/utils/ip"
	"h3s/internal/utils/logger"
)

func Delete(ctx clustercontext.ClusterContext) {
	net := network.Get(ctx)
	server, err := Get(ctx)
	if err != nil || server == nil {
		return
	}

	deleteGatewayRoute(ctx, net, server)
	delete(ctx, server)
}

func deleteGatewayRoute(ctx clustercontext.ClusterContext, net *hcloud.Network, proxy *hcloud.Server) {
	if len(proxy.PrivateNet) == 0 {
		return
	}
	ctx.Client.Network.DeleteRoute(ctx.Context, net, hcloud.NetworkDeleteRouteOpts{
		Route: hcloud.NetworkRoute{
			Destination: ip.GetIpRange("0.0.0.0/0"),
			Gateway:     proxy.PrivateNet[0].IP,
		},
	})
}

func delete(ctx clustercontext.ClusterContext, server *hcloud.Server) {
	addEvent, logEvents := logger.NewEventLogger(logger.Server, logger.Delete, server.Name)
	defer logEvents()

	_, _, err := ctx.Client.Server.DeleteWithResult(ctx.Context, server)
	if err != nil {
		addEvent(logger.Failure, err)
		return
	}

	addEvent(logger.Success)
}
