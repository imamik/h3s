package gateway

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
)

func Get(ctx *cluster.Cluster) (*hcloud.Server, error) {
	name := getName(ctx)

	l := logger.New(nil, logger.Server, logger.Get, name)
	defer l.LogEvents()

	// Get server by name
	server, _, err := ctx.CloudClient.Server.GetByName(ctx.Context, name)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	// Check if server is nil
	if server == nil {
		err = nil
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	return server, nil
}

func GetIfNeeded(ctx *cluster.Cluster) (*hcloud.Server, error) {
	if ctx.Config.PublicIps {
		return nil, nil
	}
	return Get(ctx)
}
