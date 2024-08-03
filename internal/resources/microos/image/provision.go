package image

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/clustercontext"
	"h3s/internal/resources/microos/image/commands"
	"h3s/internal/utils/logger"
	"h3s/internal/utils/ping"
	"h3s/internal/utils/ssh"
	"time"
)

const waitTime = 5 * time.Second

func Provision(
	ctx clustercontext.ClusterContext,
	architecture hcloud.Architecture,
	server *hcloud.Server,
) {
	execute(ctx, server, commands.DownloadImage(architecture), false, false)
	execute(ctx, server, commands.WriteImage(), false, true)
	execute(ctx, server, commands.Packages(), true, true)
	execute(ctx, server, commands.CleanUp(), true, false)
}

func execute(
	ctx clustercontext.ClusterContext,
	server *hcloud.Server,
	cmd string,
	pauseBefore bool,
	expectDisconnect bool,
) {
	ping.Ping(server, waitTime)

	if pauseBefore {
		logger.LogResourceEvent(
			logger.Server,
			"Execute",
			server.Name,
			logger.Initialized,
			fmt.Sprintf("Waiting for %d sec", waitTime),
			waitTime,
		)
		time.Sleep(waitTime)
	}

	_, err := ssh.ExecuteWithSsh(ctx, server, cmd)
	if err != nil {
		logger.LogResourceEvent(logger.Server, "Execute", server.Name, logger.Failure, err)
		return
	}

	if expectDisconnect {
		logger.LogResourceEvent(logger.Server, "Execute", server.Name, logger.Initialized, "Expecting disconnect")
		time.Sleep(waitTime)
	}
}
