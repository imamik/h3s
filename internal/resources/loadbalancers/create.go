package loadbalancers

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/loadbalancers/loadbalancer"
	"hcloud-k3s-cli/internal/resources/pool/node"
)

var UsePrivateIP = true

func getNodeTarget(server *hcloud.Server) hcloud.LoadBalancerCreateOptsTarget {
	return hcloud.LoadBalancerCreateOptsTarget{
		Type:         hcloud.LoadBalancerTargetTypeServer,
		Server:       hcloud.LoadBalancerCreateOptsTargetServer{Server: server},
		UsePrivateIP: &UsePrivateIP,
	}
}

func isCombined(balancerType loadbalancer.Type, n *hcloud.Server) bool {
	return balancerType == loadbalancer.Combined && (node.IsControlPlaneNode(n) || node.IsWorkerNode(n))
}

func isControlPlane(balancerType loadbalancer.Type, n *hcloud.Server) bool {
	return balancerType == loadbalancer.ControlPlane && node.IsControlPlaneNode(n)
}

func isWorker(balancerType loadbalancer.Type, n *hcloud.Server) bool {
	return balancerType == loadbalancer.Worker && node.IsWorkerNode(n)
}

func getNodeTargets(balancerType loadbalancer.Type, nodes []*hcloud.Server) []hcloud.LoadBalancerCreateOptsTarget {
	var targets []hcloud.LoadBalancerCreateOptsTarget
	for _, n := range nodes {
		if isCombined(balancerType, n) || isControlPlane(balancerType, n) || isWorker(balancerType, n) {
			targets = append(targets, getNodeTarget(n))
		}
	}
	return targets
}

func Create(ctx clustercontext.ClusterContext, network *hcloud.Network, nodes []*hcloud.Server) {
	if ctx.Config.CombinedLoadBalancer {
		create(ctx, network, nodes, loadbalancer.Combined)
	} else if ctx.Config.ControlPlane.LoadBalancer {
		create(ctx, network, nodes, loadbalancer.ControlPlane)
		create(ctx, network, nodes, loadbalancer.Worker)
	} else {
		create(ctx, network, nodes, loadbalancer.Worker)
	}
}

func create(ctx clustercontext.ClusterContext, network *hcloud.Network, nodes []*hcloud.Server, balancerType loadbalancer.Type) *hcloud.LoadBalancer {
	targets := getNodeTargets(balancerType, nodes)
	return loadbalancer.Create(ctx, network, targets, balancerType)
}
