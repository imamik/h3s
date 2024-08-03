package server

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/clustercontext"
	"h3s/internal/utils/logger"
)

const RescueModeLog = "Set Rescue Mode"

func RescueMode(
	ctx clustercontext.ClusterContext,
	sshKey *hcloud.SSHKey,
	server *hcloud.Server,
) string {
	logger.LogResourceEvent(logger.Server, RescueModeLog, server.Name, logger.Initialized)

	res, _, err := ctx.Client.Server.EnableRescue(ctx.Context, server, hcloud.ServerEnableRescueOpts{
		Type:    hcloud.ServerRescueTypeLinux64,
		SSHKeys: []*hcloud.SSHKey{sshKey},
	})
	if err != nil {
		logger.LogResourceEvent(logger.Server, RescueModeLog, server.Name, logger.Failure, err)
	}
	if err := ctx.Client.Action.WaitFor(ctx.Context, res.Action); err != nil {
		logger.LogResourceEvent(logger.Server, RescueModeLog, server.Name, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.Server, RescueModeLog, server.Name, logger.Success)

	Reboot(ctx, server)

	return res.RootPassword

}
