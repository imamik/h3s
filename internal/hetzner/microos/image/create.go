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
) (*hcloud.Image, error) {
	name := getName(ctx, architecture)
	l := logger.New(nil, logger.Image, logger.Create, name)
	l.AddEvent("...")
	l.AddEvent("This will take time")
	l.LogEvents()
	defer l.LogEvents()

	res, _, err := ctx.CloudClient.Server.CreateImage(ctx.Context, server, &hcloud.ServerCreateImageOpts{
		Type: hcloud.ImageTypeSnapshot,
		Labels: ctx.GetLabels(map[string]string{
			"is_microos":   "true",
			"architecture": string(architecture),
		}),
		Description: &name,
	})
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	if err := ctx.CloudClient.Action.WaitFor(ctx.Context, res.Action); err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	l.AddEvent(logger.Success)
	return res.Image, nil
}
