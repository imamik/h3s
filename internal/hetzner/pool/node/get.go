package node

import (
	"errors"
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/utils/logger"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Get gets a Hetzner cloud server
func Get(
	ctx *cluster.Cluster,
	pool config.NodePool,
	i int,
) (*hcloud.Server, error) {
	name := getName(ctx, pool, i)

	l := logger.New(nil, logger.Server, logger.Get, name)
	defer l.LogEvents()

	server, _, err := ctx.CloudClient.Server.GetByName(ctx.Context, name)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}
	if server == nil {
		err = errors.New("server is nil")
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	l.AddEvent(logger.Success)
	return server, nil
}
