package image

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/microos/image/commands"
	"hcloud-k3s-cli/internal/utils/logger"
	"hcloud-k3s-cli/internal/utils/ping"
	"time"
)

func Provision(
	ctx clustercontext.ClusterContext,
	architecture hcloud.Architecture,
	server *hcloud.Server,
) {
	execute(ctx, server, commands.DownloadImage(architecture), 0, false)
	execute(ctx, server, commands.WriteImage(), 0, true)
	execute(ctx, server, commands.Packages(), 5, true)
	execute(ctx, server, commands.CleanUp(), 5, false)
}

func execute(
	ctx clustercontext.ClusterContext,
	server *hcloud.Server,
	cmd string,
	pauseBeforeSeconds int,
	expectDisconnect bool,
) {
	if pauseBeforeSeconds > 0 {
		logger.LogResourceEvent(logger.Server, "Execute", server.Name, logger.Initialized, "Waiting for %d sec", pauseBeforeSeconds)
		time.Sleep(time.Duration(pauseBeforeSeconds) * time.Second)
	}

	if expectDisconnect {
		logger.LogResourceEvent(logger.Server, "Execute", server.Name, logger.Initialized, "Expecting disconnect")
		time.Sleep(5 * time.Second)
		ping.Ping(server)
	}
}
