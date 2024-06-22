package placementgroup

import (
	"hcloud-k3s-cli/pkg/config"
	"hcloud-k3s-cli/pkg/resources/clustercontext"
	"hcloud-k3s-cli/pkg/utils/logger"
)

func Delete(
	ctx clustercontext.ClusterContext,
	pool config.NodePool,
) {
	placementGroup := Get(ctx, pool)
	if placementGroup == nil {
		return
	}

	logger.LogResourceEvent(logger.PlacementGroup, logger.Delete, placementGroup.Name, logger.Initialized)

	_, err := ctx.Client.PlacementGroup.Delete(ctx.Context, placementGroup)
	if err != nil {
		logger.LogResourceEvent(logger.PlacementGroup, logger.Delete, placementGroup.Name, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.PlacementGroup, logger.Delete, placementGroup.Name, logger.Success)
}
