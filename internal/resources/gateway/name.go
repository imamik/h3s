package gateway

import "hcloud-k3s-cli/internal/clustercontext"

func getName(ctx clustercontext.ClusterContext) string {
	return ctx.GetName("gateway")
}
