package placementgroup

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/config"
	"hcloud-k3s-cli/pkg/resources/clustercontext"
	"hcloud-k3s-cli/pkg/utils/logger"
)

func Get(
	ctx clustercontext.ClusterContext,
	pool config.NodePool,
) *hcloud.PlacementGroup {
	name := getName(ctx, pool)
	logger.LogResourceEvent(logger.PlacementGroup, logger.Get, name, logger.Initialized)

	placementGroup, _, err := ctx.Client.PlacementGroup.GetByName(ctx.Context, name)
	if err != nil || placementGroup == nil {
		logger.LogResourceEvent(logger.PlacementGroup, logger.Get, name, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.PlacementGroup, logger.Get, name, logger.Success)
	return placementGroup
}
