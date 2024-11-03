package placementgroup

import (
	"errors"
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/utils/logger"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Get gets the Hetzner cloud placement group
func Get(
	ctx *cluster.Cluster,
	pool config.NodePool,
) (*hcloud.PlacementGroup, error) {
	name := getName(ctx, pool)

	l := logger.New(nil, logger.PlacementGroup, logger.Get, name)
	defer l.LogEvents()

	placementGroup, _, err := ctx.CloudClient.PlacementGroup.GetByName(ctx.Context, name)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}
	if placementGroup == nil {
		err = errors.New("placement group is nil")
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	l.AddEvent(logger.Success)
	return placementGroup, nil
}
