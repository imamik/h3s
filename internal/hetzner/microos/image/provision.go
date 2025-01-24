package image

import (
	"fmt"
	"h3s/internal/cluster"
	"h3s/internal/hetzner/microos/image/commands"
	"h3s/internal/utils/logger"
	"h3s/internal/utils/ping"
	"h3s/internal/utils/ssh"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

const waitTime = 5 * time.Second

// Provision provisions the Hetzner cloud microOS image
func Provision(
	ctx *cluster.Cluster,
	architecture hcloud.Architecture,
	server *hcloud.Server,
) error {
	if err := execute(ctx, server, commands.DownloadImage(architecture), false, false); err != nil {
		return err
	}
	if err := execute(ctx, server, commands.WriteImage(), false, true); err != nil {
		return err
	}
	if err := execute(ctx, server, commands.Packages(), true, true); err != nil {
		return err
	}
	if err := execute(ctx, server, commands.CleanUp(), true, false); err != nil {
		return err
	}
	return nil
}

func execute(
	ctx *cluster.Cluster,
	server *hcloud.Server,
	cmd string,
	pauseBefore bool,
	expectDisconnect bool,
) error {
	l := logger.New(nil, logger.Server, "Execute", server.Name)
	ping.Ping(server, waitTime)

	if pauseBefore {
		l.AddEvent(logger.Initialized, fmt.Sprintf("Waiting for %d sec", waitTime))
		time.Sleep(waitTime)
	}

	_, err := ssh.ExecuteWithSSH(ctx.Config.SSHKeyPaths.PrivateKeyPath, server, cmd)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	if expectDisconnect {
		l.AddEvent(logger.Initialized, "Expecting disconnect")
		time.Sleep(waitTime)
		l.AddEvent(logger.Success, "Added wait time")
	}

	return nil
}
