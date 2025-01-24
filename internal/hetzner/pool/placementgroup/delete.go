package placementgroup

import (
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/utils/logger"
)

// Delete deletes a Hetzner cloud placement group
func Delete(
	ctx *cluster.Cluster,
	pool config.NodePool,
) error {
	l := logger.New(nil, logger.PlacementGroup, logger.Delete, ctx.GetName(pool.Name))
	defer l.LogEvents()

	placementGroup, err := Get(ctx, pool)
	if placementGroup == nil && err.Error() == "placement group is nil" {
		l.AddEvent(logger.Success, "no placement group found to delete")
		return nil
	}

	_, err = ctx.CloudClient.PlacementGroup.Delete(ctx.Context, placementGroup)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	l.AddEvent(logger.Success)
	return nil
}
