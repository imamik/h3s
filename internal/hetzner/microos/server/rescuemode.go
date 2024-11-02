package server

import (
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

const RescueModeLog = "Set Rescue Mode"

func RescueMode(
	ctx *cluster.Cluster,
	sshKey *hcloud.SSHKey,
	server *hcloud.Server,
) (string, error) {
	l := logger.New(nil, logger.Server, RescueModeLog, server.Name)
	defer l.LogEvents()

	res, _, err := ctx.CloudClient.Server.EnableRescue(ctx.Context, server, hcloud.ServerEnableRescueOpts{
		Type:    hcloud.ServerRescueTypeLinux64,
		SSHKeys: []*hcloud.SSHKey{sshKey},
	})
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return "", err
	}
	if err := ctx.CloudClient.Action.WaitFor(ctx.Context, res.Action); err != nil {
		l.AddEvent(logger.Failure, err)
		return "", err
	}

	l.AddEvent(logger.Success)
	if err := Reboot(ctx, server); err != nil {
		return "", err
	}

	return res.RootPassword, nil
}
