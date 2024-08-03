package server

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/clustercontext"
	"h3s/internal/utils/logger"
)

const RebootLog = "Reboot"

func Reboot(ctx clustercontext.ClusterContext, server *hcloud.Server) {
	logger.LogResourceEvent(logger.Server, RebootLog, server.Name, logger.Initialized)

	action, _, err := ctx.Client.Server.Reset(ctx.Context, server)
	if err != nil {
		logger.LogResourceEvent(logger.Server, RebootLog, server.Name, logger.Failure, err)
	}
	if err := ctx.Client.Action.WaitFor(ctx.Context, action); err != nil {
		logger.LogResourceEvent(logger.Server, RebootLog, server.Name, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.Server, RebootLog, server.Name, logger.Success)
}
