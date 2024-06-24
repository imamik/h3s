package loadbalancer

import (
	"hcloud-k3s-cli/internal/clustercontext"
)

func getName(ctx clustercontext.ClusterContext, balancerType Type) string {
	return ctx.GetName(string(balancerType), "loadbalancer")
}
