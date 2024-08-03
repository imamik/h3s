package placementgroup

import (
	"h3s/internal/clustercontext"
	"h3s/internal/config"
)

func getName(
	ctx clustercontext.ClusterContext,
	pool config.NodePool,
) string {
	return ctx.GetName(pool.Name, "pool")
}
