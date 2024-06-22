package placementgroup

import (
	"hcloud-k3s-cli/pkg/config"
	"hcloud-k3s-cli/pkg/resources/clustercontext"
)

func getName(
	ctx clustercontext.ClusterContext,
	pool config.NodePool,
) string {
	return ctx.GetName(pool.Name, "pool")
}
