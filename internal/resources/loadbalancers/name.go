package loadbalancers

import (
	"h3s/internal/clustercontext"
)

func getName(ctx clustercontext.ClusterContext) string {
	return ctx.GetName("loadbalancer")
}
