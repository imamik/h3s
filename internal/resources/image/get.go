package image

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Get(ctx clustercontext.ClusterContext, architecture hcloud.Architecture) (*hcloud.Image, error) {
	name := getName(ctx, architecture)
	logger.LogResourceEvent(logger.Image, logger.Get, name, logger.Initialized)

	snapshot, _, err := ctx.Client.Image.GetForArchitecture(ctx.Context, name, architecture)
	if err != nil {
		logger.LogResourceEvent(logger.Image, logger.Get, name, logger.Failure, err)
		return nil, err
	}

	logger.LogResourceEvent(logger.Image, logger.Get, name, logger.Success)
	return snapshot, nil
}
