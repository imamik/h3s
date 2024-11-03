// Package placementgroup contains the functionality for managing Hetzner cloud placement groups
package placementgroup

import (
	"errors"
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/utils/logger"
	"strconv"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func create(
	ctx *cluster.Cluster,
	pool config.NodePool,
	isControlPlane bool,
	isWorker bool,
) (*hcloud.PlacementGroup, error) {
	name := getName(ctx, pool)

	l := logger.New(nil, logger.PlacementGroup, logger.Create, name)
	defer l.LogEvents()

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
		l.AddEvent(logger.Failure, err)
		return nil, err
	}
	if res.PlacementGroup == nil {
		err = errors.New("placement group is nil")
		l.AddEvent(logger.Failure, err)
		return nil, err
	}
	if err := ctx.CloudClient.Action.WaitFor(ctx.Context, res.Action); err != nil {
		l.AddEvent(logger.Failure, err)
	}

	l.AddEvent(logger.Success)
	return res.PlacementGroup, nil
}

// Create creates a Hetzner cloud placement group
func Create(
	ctx *cluster.Cluster,
	pool config.NodePool,
	isControlPlane bool,
	isWorker bool,
) (*hcloud.PlacementGroup, error) {
	placementGroup, err := Get(ctx, pool)
	if placementGroup != nil && err == nil {
		return placementGroup, nil
	}
	return create(ctx, pool, isControlPlane, isWorker)
}
