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
	addEvent, logEvents := logger.NewEventLogger(logger.PlacementGroup, logger.Create, name)
	defer logEvents()

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
		addEvent(logger.Failure, err)
	}
	if res.PlacementGroup == nil {
		addEvent(logger.Failure, "Empty Response")
	}
	if err := ctx.Client.Action.WaitFor(ctx.Context, res.Action); err != nil {
		addEvent(logger.Failure, err)
	}

	addEvent(logger.Success)
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
