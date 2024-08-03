package node

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/clustercontext"
	"h3s/internal/config"
	"h3s/internal/utils/logger"
)

func Get(
	ctx clustercontext.ClusterContext,
	pool config.NodePool,
	i int,
) *hcloud.Server {
	name := getName(ctx, pool, i)
	addEvent, logEvents := logger.NewEventLogger(logger.Server, logger.Get, name)
	defer logEvents()

	server, _, err := ctx.Client.Server.GetByName(ctx.Context, name)
	if err != nil || server == nil {
		addEvent(logger.Failure, err)
		return nil
	}

	addEvent(logger.Success)
	return server
}
