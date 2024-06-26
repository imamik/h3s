package image

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/image/server"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Delete(ctx clustercontext.ClusterContext, architecture hcloud.Architecture) {
	deleteImage(ctx, architecture)
	server.Delete(ctx, architecture)
}

func deleteImage(ctx clustercontext.ClusterContext, architecture hcloud.Architecture) {
	img, err := Get(ctx, architecture)
	if err != nil || img == nil {
		return
	}

	logger.LogResourceEvent(logger.Image, logger.Delete, img.Name, logger.Initialized)

	_, err = ctx.Client.Image.Delete(ctx.Context, img)
	if err != nil {
		logger.LogResourceEvent(logger.Image, logger.Delete, img.Name, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.Image, logger.Delete, img.Name, logger.Success)
}
