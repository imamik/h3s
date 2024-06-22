package placementgroup

import (
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/config"
)

func getName(
	ctx clustercontext.ClusterContext,
	pool config.NodePool,
) string {
	return ctx.GetName(pool.Name, "pool")
}
