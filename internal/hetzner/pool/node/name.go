package node

import (
	"h3s/internal/cluster"
	"h3s/internal/config"
	"strconv"
)

func getName(
	ctx *cluster.Cluster,
	pool config.NodePool,
	index int,
) string {
	indexStr := strconv.Itoa(index)
	return ctx.GetName(pool.Name, "node", indexStr)
}
