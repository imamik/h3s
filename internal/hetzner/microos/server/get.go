package server

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
)

func Get(
	ctx *cluster.Cluster,
	architecture hcloud.Architecture,
) *hcloud.Server {
	serverName := getName(ctx, architecture)
	logger.LogResourceEvent(logger.Server, logger.Get, serverName, logger.Initialized)

	server, _, err := ctx.CloudClient.Server.GetByName(ctx.Context, serverName)
	if err != nil || server == nil {
		logger.LogResourceEvent(logger.Server, logger.Get, serverName, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.Server, logger.Get, serverName, logger.Success)
	return server
}
