package loadbalancer

import (
	"hcloud-k3s-cli/pkg/clustercontext"
)

func getName(ctx clustercontext.ClusterContext, balancerType Type) string {
	return ctx.GetName(string(balancerType), "loadbalancer")
}
