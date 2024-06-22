package placementgroup

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/config"
	"hcloud-k3s-cli/pkg/utils/logger"
)

func create(
	ctx clustercontext.ClusterContext,
	pool config.NodePool,
) *hcloud.PlacementGroup {
	placementGroupName := getName(ctx, pool)
	logger.LogResourceEvent(logger.PlacementGroup, logger.Create, placementGroupName, logger.Initialized)

	placementGroupResp, _, err := ctx.Client.PlacementGroup.Create(ctx.Context, hcloud.PlacementGroupCreateOpts{
		Name:   placementGroupName,
		Type:   hcloud.PlacementGroupTypeSpread,
		Labels: ctx.GetLabels(),
	})
	if err != nil {
		logger.LogResourceEvent(logger.PlacementGroup, logger.Create, placementGroupName, logger.Failure, err)
	}
	if placementGroupResp.PlacementGroup == nil {
		logger.LogResourceEvent(logger.PlacementGroup, logger.Create, placementGroupName, logger.Failure, "Empty Response")
	}

	logger.LogResourceEvent(logger.PlacementGroup, logger.Create, placementGroupName, logger.Success)
	return placementGroupResp.PlacementGroup
}

func Create(
	ctx clustercontext.ClusterContext,
	pool config.NodePool,
) *hcloud.PlacementGroup {
	placementGroup := Get(ctx, pool)
	if placementGroup == nil {
		placementGroup = create(ctx, pool)
	}
	return placementGroup
}
