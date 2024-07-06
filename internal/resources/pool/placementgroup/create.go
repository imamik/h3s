package placementgroup

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/config"
	"hcloud-k3s-cli/internal/utils/logger"
	"strconv"
)

func create(
	ctx clustercontext.ClusterContext,
	pool config.NodePool,
	isControlPlane bool,
	isWorker bool,
) *hcloud.PlacementGroup {
	name := getName(ctx, pool)
	logger.LogResourceEvent(logger.PlacementGroup, logger.Create, name, logger.Initialized)

	res, _, err := ctx.Client.PlacementGroup.Create(ctx.Context, hcloud.PlacementGroupCreateOpts{
		Name: name,
		Type: hcloud.PlacementGroupTypeSpread,
		Labels: ctx.GetLabels(map[string]string{
			"pool":             pool.Name,
			"nodes":            strconv.Itoa(pool.Nodes),
			"is_control_plane": strconv.FormatBool(isControlPlane),
			"is_worker":        strconv.FormatBool(isWorker),
		}),
	})
	if err != nil {
		logger.LogResourceEvent(logger.PlacementGroup, logger.Create, name, logger.Failure, err)
	}
	if res.PlacementGroup == nil {
		logger.LogResourceEvent(logger.PlacementGroup, logger.Create, name, logger.Failure, "Empty Response")
	}
	if err := ctx.Client.Action.WaitFor(ctx.Context, res.Action); err != nil {
		logger.LogResourceEvent(logger.PlacementGroup, logger.Create, name, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.PlacementGroup, logger.Create, name, logger.Success)
	return res.PlacementGroup
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
