package node

import (
	"h3s/internal/clustercontext"
	"h3s/internal/config"
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
