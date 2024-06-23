package placementgroup

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/clustercontext"
	"hcloud-k3s-cli/pkg/config"
	"hcloud-k3s-cli/pkg/utils/logger"
	"strconv"
)

func create(
	ctx clustercontext.ClusterContext,
	pool config.NodePool,
	isControlPlane bool,
	isWorker bool,
) *hcloud.PlacementGroup {
	placementGroupName := getName(ctx, pool)
	logger.LogResourceEvent(logger.PlacementGroup, logger.Create, placementGroupName, logger.Initialized)

	placementGroupResp, _, err := ctx.Client.PlacementGroup.Create(ctx.Context, hcloud.PlacementGroupCreateOpts{
		Name: placementGroupName,
		Type: hcloud.PlacementGroupTypeSpread,
		Labels: ctx.GetLabels(map[string]string{
			"pool":             pool.Name,
			"nodes":            strconv.Itoa(pool.Nodes),
			"is_control_plane": strconv.FormatBool(isControlPlane),
			"is_worker":        strconv.FormatBool(isWorker),
		}),
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
	isControlPlane bool,
	isWorker bool,
) *hcloud.PlacementGroup {
	placementGroup := Get(ctx, pool)
	if placementGroup == nil {
		placementGroup = create(ctx, pool, isControlPlane, isWorker)
	}
	return placementGroup
}
