package image

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
)

func Delete(ctx *cluster.Cluster, architecture hcloud.Architecture) error {
	l := logger.New(nil, logger.Image, logger.Delete, getName(ctx, architecture))
	defer l.LogEvents()

	img, err := Get(ctx, architecture)

	if img == nil && err == nil {
		l.AddEvent(logger.Success, "no image found to delete")
		return nil
	}

	if err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	_, err = ctx.CloudClient.Image.Delete(ctx.Context, img)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	l.AddEvent(logger.Success)
	return nil
}
