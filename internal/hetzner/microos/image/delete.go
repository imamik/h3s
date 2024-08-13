package image

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
)

func Delete(ctx *cluster.Cluster, architecture hcloud.Architecture) {
	img, err := Get(ctx, architecture)
	if err != nil || img == nil {
		return
	}

	logger.LogResourceEvent(logger.Image, logger.Delete, img.Name, logger.Initialized)

	_, err = ctx.CloudClient.Image.Delete(ctx.Context, img)
	if err != nil {
		logger.LogResourceEvent(logger.Image, logger.Delete, img.Name, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.Image, logger.Delete, img.Name, logger.Success)
}
