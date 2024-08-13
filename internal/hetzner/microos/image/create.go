package image

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
)

func Create(
	ctx *cluster.Cluster,
	server *hcloud.Server,
	architecture hcloud.Architecture,
) *hcloud.Image {
	name := getName(ctx, architecture)
	logger.LogResourceEvent(logger.Image, logger.Create, name, logger.Initialized)
	logger.LogResourceEvent(logger.Image, "...", name, logger.Initialized)
	logger.LogResourceEvent(logger.Image, "This will take time", name, logger.Success)

	res, _, err := ctx.CloudClient.Server.CreateImage(ctx.Context, server, &hcloud.ServerCreateImageOpts{
		Type: hcloud.ImageTypeSnapshot,
		Labels: ctx.GetLabels(map[string]string{
			"is_microos":   "true",
			"architecture": string(architecture),
		}),
		Description: &name,
	})
	if err != nil {
		logger.LogResourceEvent(logger.Image, logger.Create, name, logger.Failure, err)
	}

	if err := ctx.CloudClient.Action.WaitFor(ctx.Context, res.Action); err != nil {
		logger.LogResourceEvent(logger.Image, logger.Create, name, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.Image, logger.Create, name, logger.Success)

	return res.Image
}
