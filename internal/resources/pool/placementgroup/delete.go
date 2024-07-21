package placementgroup

import (
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/config"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Delete(
	ctx clustercontext.ClusterContext,
	pool config.NodePool,
) {
	addEvent, logEvents := logger.NewEventLogger(logger.PlacementGroup, logger.Delete, ctx.GetName(pool.Name))
	defer logEvents()

	placementGroup := Get(ctx, pool)
	if placementGroup == nil {
		return
	}

	_, err := ctx.Client.PlacementGroup.Delete(ctx.Context, placementGroup)
	if err != nil {
		addEvent(logger.Failure, err)
		return
	}

	addEvent(logger.Success)
}
