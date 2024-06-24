package placementgroup

import (
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/config"
)

func getName(
	ctx clustercontext.ClusterContext,
	pool config.NodePool,
) string {
	return ctx.GetName(pool.Name, "pool")
}
