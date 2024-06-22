package server

import (
	"hcloud-k3s-cli/pkg/config"
	"hcloud-k3s-cli/pkg/resources/clustercontext"
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
