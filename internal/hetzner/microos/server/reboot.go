package server

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
)

const RebootLog = "Reboot"

func Reboot(ctx *cluster.Cluster, server *hcloud.Server) error {
	l := logger.New(nil, logger.Server, RebootLog, server.Name)
	defer l.LogEvents()

	action, _, err := ctx.CloudClient.Server.Reset(ctx.Context, server)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}
	if err := ctx.CloudClient.Action.WaitFor(ctx.Context, action); err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	l.AddEvent(logger.Success)
	return nil
}
