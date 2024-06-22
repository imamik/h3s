package server

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/config"
	"hcloud-k3s-cli/pkg/utils/logger"
)

func Get(
	ctx clustercontext.ClusterContext,
	pool config.NodePool,
	i int,
) *hcloud.Server {
	serverName := getName(ctx, pool, i)
	logger.LogResourceEvent(logger.Server, logger.Get, serverName, logger.Initialized)

	server, _, err := ctx.Client.Server.GetByName(ctx.Context, serverName)
	if err != nil || server == nil {
		logger.LogResourceEvent(logger.Server, logger.Get, serverName, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.Server, logger.Get, serverName, logger.Success)
	return server
}
