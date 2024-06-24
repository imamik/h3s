package node

import (
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/config"
	"strconv"
)

func getName(
	ctx clustercontext.ClusterContext,
	pool config.NodePool,
	index int,
) string {
	indexStr := strconv.Itoa(index)
	return ctx.GetName(pool.Name, "node", indexStr)
}
