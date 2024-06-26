package server

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Get(
	ctx clustercontext.ClusterContext,
	architecture hcloud.Architecture,
) *hcloud.Server {
	serverName := getName(ctx, architecture)
	logger.LogResourceEvent(logger.Server, logger.Get, serverName, logger.Initialized)

	server, _, err := ctx.Client.Server.GetByName(ctx.Context, serverName)
	if err != nil || server == nil {
		logger.LogResourceEvent(logger.Server, logger.Get, serverName, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.Server, logger.Get, serverName, logger.Success)
	return server
}
