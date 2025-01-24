package placementgroup

import (
	"h3s/internal/cluster"
	"h3s/internal/config"
)

func getName(
	ctx *cluster.Cluster,
	pool config.NodePool,
) string {
	return ctx.GetName(pool.Name, "pool")
}
