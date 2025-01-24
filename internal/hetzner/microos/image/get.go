package image

import (
	"errors"
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Get gets the Hetzner cloud microOS image
func Get(ctx *cluster.Cluster, architecture hcloud.Architecture) (*hcloud.Image, error) {
	name := getName(ctx, architecture)

	l := logger.New(nil, logger.Image, logger.Get, name)
	defer l.LogEvents()

	options := hcloud.ImageListOpts{
		Type:         []hcloud.ImageType{hcloud.ImageTypeSnapshot},
		Architecture: []hcloud.Architecture{architecture},
	}

	images, err := ctx.CloudClient.Image.AllWithOpts(ctx.Context, options)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}
	l.AddEvent(logger.Success, "found images")

	// Find the correct image
	var image *hcloud.Image
	description := getName(ctx, architecture)
	for _, img := range images {
		if img.Description == description {
			image = img
			break
		}
	}

	if image == nil {
		err := errors.New("image is nil")
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	return image, nil
}
