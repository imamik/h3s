package server

import (
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

const ShutDownLog = "ShutDown"

func Shutdown(ctx *cluster.Cluster, server *hcloud.Server) error {
	l := logger.New(nil, logger.Server, ShutDownLog, server.Name)
	defer l.LogEvents()

	action, _, err := ctx.CloudClient.Server.Shutdown(ctx.Context, server)
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
