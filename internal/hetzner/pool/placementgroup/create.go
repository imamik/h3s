package placementgroup

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/utils/logger"
	"strconv"
)

func create(
	ctx *cluster.Cluster,
	pool config.NodePool,
	isControlPlane bool,
	isWorker bool,
) *hcloud.PlacementGroup {
	name := getName(ctx, pool)
	addEvent, logEvents := logger.NewEventLogger(logger.PlacementGroup, logger.Create, name)
	defer logEvents()

	res, _, err := ctx.CloudClient.PlacementGroup.Create(ctx.Context, hcloud.PlacementGroupCreateOpts{
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
	if err := ctx.CloudClient.Action.WaitFor(ctx.Context, res.Action); err != nil {
		addEvent(logger.Failure, err)
	}

	addEvent(logger.Success)
	return res.PlacementGroup
}

func Create(
	ctx *cluster.Cluster,
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
