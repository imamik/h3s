package image

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
)

func Get(ctx *cluster.Cluster, architecture hcloud.Architecture) (*hcloud.Image, error) {
	name := getName(ctx, architecture)
	logger.LogResourceEvent(logger.Image, logger.Get, name, logger.Initialized)

	options := hcloud.ImageListOpts{
		Type:         []hcloud.ImageType{hcloud.ImageTypeSnapshot},
		Architecture: []hcloud.Architecture{architecture},
	}

	images, err := ctx.CloudClient.Image.AllWithOpts(ctx.Context, options)

	if err != nil {
		logger.LogResourceEvent(logger.Image, logger.Get, name, logger.Failure, err)
		return nil, err
	}

	logger.LogResourceEvent(logger.Image, logger.Get, name, logger.Success)

	// Find the correct image
	var image *hcloud.Image
	description := getName(ctx, architecture)
	for _, img := range images {
		if img.Description == description {
			image = img
			break
		}
	}
	return image, nil
}
