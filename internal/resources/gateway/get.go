package gateway

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Get(ctx clustercontext.ClusterContext) (*hcloud.Server, error) {
	name := getName(ctx)
	logger.LogResourceEvent(logger.Server, logger.Get, name, logger.Initialized)

	server, _, err := ctx.Client.Server.GetByName(ctx.Context, name)
	if err != nil {
		logger.LogResourceEvent(logger.Server, logger.Get, name, logger.Failure, err)
		return nil, err
	}

	logger.LogResourceEvent(logger.Server, logger.Get, name, logger.Success)
	return server, nil
}
