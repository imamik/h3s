package placementgroup

import (
	"h3s/internal/clustercontext"
	"h3s/internal/config"
	"h3s/internal/utils/logger"
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
