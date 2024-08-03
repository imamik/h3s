package gateway

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/clustercontext"
	"h3s/internal/utils/logger"
)

func Get(ctx clustercontext.ClusterContext) (*hcloud.Server, error) {
	name := getName(ctx)
	addEvent, logEvents := logger.NewEventLogger(logger.Server, logger.Create, name)
	defer logEvents()

	server, _, err := ctx.Client.Server.GetByName(ctx.Context, name)
	if err != nil {
		addEvent(logger.Failure, err)
		return nil, err
	}

	addEvent(logger.Success)
	return server, nil
}
