package placementgroup

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/utils/logger"
)

func Get(
	ctx *cluster.Cluster,
	pool config.NodePool,
) *hcloud.PlacementGroup {
	name := getName(ctx, pool)
	addEvent, logEvents := logger.NewEventLogger(logger.PlacementGroup, logger.Get, name)
	defer logEvents()

	placementGroup, _, err := ctx.CloudClient.PlacementGroup.GetByName(ctx.Context, name)
	if err != nil || placementGroup == nil {
		addEvent(logger.Failure, err)
		return nil
	}

	addEvent(logger.Success)
	return placementGroup
}
