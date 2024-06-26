package image

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Create(
	ctx clustercontext.ClusterContext,
	server *hcloud.Server,
	architecture hcloud.Architecture,
) *hcloud.Image {
	name := getName(ctx, architecture)
	logger.LogResourceEvent(logger.Image, logger.Create, name, logger.Initialized)
	logger.LogResourceEvent(logger.Image, "...", name, logger.Initialized)
	logger.LogResourceEvent(logger.Image, "This will take time", name, logger.Success)

	res, _, err := ctx.Client.Server.CreateImage(ctx.Context, server, &hcloud.ServerCreateImageOpts{
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

	if err := ctx.Client.Action.WaitFor(ctx.Context, res.Action); err != nil {
		logger.LogResourceEvent(logger.Image, logger.Create, name, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.Image, logger.Create, name, logger.Success)

	return res.Image
}
