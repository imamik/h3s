package server

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
)

const RebootLog = "Reboot"

func Reboot(ctx *cluster.Cluster, server *hcloud.Server) {
	logger.LogResourceEvent(logger.Server, RebootLog, server.Name, logger.Initialized)

	action, _, err := ctx.CloudClient.Server.Reset(ctx.Context, server)
	if err != nil {
		logger.LogResourceEvent(logger.Server, RebootLog, server.Name, logger.Failure, err)
	}
	if err := ctx.CloudClient.Action.WaitFor(ctx.Context, action); err != nil {
		logger.LogResourceEvent(logger.Server, RebootLog, server.Name, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.Server, RebootLog, server.Name, logger.Success)
}
