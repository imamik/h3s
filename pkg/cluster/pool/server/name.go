package server

import (
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/config"
)

func getName(pool config.NodePool, index int, ctx clustercontext.ClusterContext) string {
	return ctx.GetName(pool.Name, "node", string(rune(index)))
}
