package gateway

import "h3s/internal/cluster"

func getName(ctx *cluster.Cluster) string {
	return ctx.GetName("gateway")
}
