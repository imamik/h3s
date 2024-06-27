package image

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/resources/microos/image/commands"
	"hcloud-k3s-cli/internal/utils/logger"
	"hcloud-k3s-cli/internal/utils/ping"
	"hcloud-k3s-cli/internal/utils/ssh"
	"time"
)

func Provision(
	architecture hcloud.Architecture,
	server *hcloud.Server,
	rootPassword string,
) {
	execute(server, rootPassword, commands.DownloadImage(architecture), 0, false)
	execute(server, rootPassword, commands.WriteImage(), 0, true)
	execute(server, rootPassword, commands.Packages(), 5, true)
	execute(server, rootPassword, commands.CleanUp(), 5, false)
}

func execute(
	server *hcloud.Server,
	rootPassword string,
	cmd string,
	pauseBeforeSeconds time.Duration,
	expectDisconnect bool,
) {
	ping.Ping(server, 5)

	if pauseBeforeSeconds > 0 {
		logger.LogResourceEvent(
			logger.Server,
			"Execute",
			server.Name,
			logger.Initialized,
			fmt.Sprintf("Waiting for %d sec", pauseBeforeSeconds),
			pauseBeforeSeconds,
		)
		time.Sleep(pauseBeforeSeconds * time.Second)
	}

	_, err := ssh.ExecuteWithPassword(server, rootPassword, cmd)
	if err != nil {
		logger.LogResourceEvent(logger.Server, "Execute", server.Name, logger.Failure, err)
		return
	}

	if expectDisconnect {
		logger.LogResourceEvent(logger.Server, "Execute", server.Name, logger.Initialized, "Expecting disconnect")
		time.Sleep(5 * time.Second)
		ping.Ping(server, 5)
	}
}
