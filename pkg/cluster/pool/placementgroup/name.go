package placementgroup

import (
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/config"
)

func getName(pool config.NodePool, ctx clustercontext.ClusterContext) string {
	return ctx.GetName(pool.Name, "pool")
}
