package gateway

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
)

func Get(ctx *cluster.Cluster) (*hcloud.Server, error) {
	name := getName(ctx)
	addEvent, logEvents := logger.NewEventLogger(logger.Server, logger.Create, name)
	defer logEvents()

	server, _, err := ctx.CloudClient.Server.GetByName(ctx.Context, name)
	if err != nil {
		addEvent(logger.Failure, err)
		return nil, err
	}

	addEvent(logger.Success)
	return server, nil
}

func GetIfNeeded(ctx *cluster.Cluster) (*hcloud.Server, error) {
	if ctx.Config.PublicIps {
		return nil, nil
	}
	return Get(ctx)
}
