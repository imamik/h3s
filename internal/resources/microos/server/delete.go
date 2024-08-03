package server

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/clustercontext"
	"h3s/internal/utils/logger"
)

func Delete(
	ctx clustercontext.ClusterContext,
	architecture hcloud.Architecture,
) {
	server := Get(ctx, architecture)

	if server == nil {
		return
	}

	logger.LogResourceEvent(logger.Server, logger.Delete, server.Name, logger.Initialized)

	_, _, err := ctx.Client.Server.DeleteWithResult(ctx.Context, server)
	if err != nil {
		logger.LogResourceEvent(logger.Server, logger.Delete, server.Name, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.Server, logger.Delete, server.Name, logger.Success)
}
