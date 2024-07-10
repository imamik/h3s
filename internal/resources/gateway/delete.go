package gateway

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/network"
	"hcloud-k3s-cli/internal/utils/ip"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Delete(ctx clustercontext.ClusterContext) {
	net := network.Get(ctx)
	server, err := getServer(ctx)
	if err != nil || server == nil {
		return
	}

	deleteGatewayRoute(ctx, net, server)
	delete(ctx, server)
}

func deleteGatewayRoute(ctx clustercontext.ClusterContext, net *hcloud.Network, proxy *hcloud.Server) {
	ctx.Client.Network.DeleteRoute(ctx.Context, net, hcloud.NetworkDeleteRouteOpts{
		Route: hcloud.NetworkRoute{
			Destination: ip.GetIpRange("0.0.0.0/0"),
			Gateway:     proxy.PrivateNet[0].IP,
		},
	})
}

func delete(ctx clustercontext.ClusterContext, server *hcloud.Server) {
	logger.LogResourceEvent(logger.Server, logger.Delete, server.Name, logger.Initialized)

	_, _, err := ctx.Client.Server.DeleteWithResult(ctx.Context, server)
	if err != nil {
		logger.LogResourceEvent(logger.Server, logger.Delete, server.Name, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.Server, logger.Delete, server.Name, logger.Success)
}
