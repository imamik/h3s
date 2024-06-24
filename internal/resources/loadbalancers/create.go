package loadbalancers

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/loadbalancers/loadbalancer"
	"hcloud-k3s-cli/internal/resources/pool/node"
)

var UsePrivateIP = true

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

func nodeInLoadBalancer(lb *hcloud.LoadBalancer, node *hcloud.Server) bool {
	for _, t := range lb.Targets {
		if t.Server.Server.ID == node.ID {
			return true
		}
	}
	return false
}

func addTarget(ctx clustercontext.ClusterContext, lb *hcloud.LoadBalancer, node *hcloud.Server) {
	if nodeInLoadBalancer(lb, node) {
		return
	}
	_, _, _ = ctx.Client.LoadBalancer.AddServerTarget(ctx.Context, lb, hcloud.LoadBalancerAddServerTargetOpts{
		Server:       node,
		UsePrivateIP: &UsePrivateIP,
	})
}

func create(ctx clustercontext.ClusterContext, network *hcloud.Network, nodes []*hcloud.Server, balancerType loadbalancer.Type) {
	lb := loadbalancer.Create(ctx, network, balancerType)
	for _, n := range nodes {
		if balancerType == loadbalancer.Combined {
			addTarget(ctx, lb, n)
		} else if balancerType == loadbalancer.ControlPlane && node.IsControlPlaneNode(n) {
			addTarget(ctx, lb, n)
		} else if balancerType == loadbalancer.Worker && node.IsWorkerNode(n) {
			addTarget(ctx, lb, n)
		}
	}
}
