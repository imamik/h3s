package server

import (
	"errors"
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Get gets the Hetzner cloud microOS server
func Get(
	ctx *cluster.Cluster,
	architecture hcloud.Architecture,
) (*hcloud.Server, error) {
	serverName := getName(ctx, architecture)
	l := logger.New(nil, logger.Server, logger.Get, serverName)
	defer l.LogEvents()

	server, _, err := ctx.CloudClient.Server.GetByName(ctx.Context, serverName)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}
	if server == nil {
		err := errors.New("server is nil")
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	l.AddEvent(logger.Success)
	return server, nil
}
