package loadbalancers

import (
	"hcloud-k3s-cli/pkg/clustercontext"
	"hcloud-k3s-cli/pkg/resources/loadbalancers/loadbalancer"
)

func Delete(ctx clustercontext.ClusterContext) {
	loadbalancer.Delete(ctx, loadbalancer.Combined)
	loadbalancer.Delete(ctx, loadbalancer.ControlPlane)
	loadbalancer.Delete(ctx, loadbalancer.Worker)
}
