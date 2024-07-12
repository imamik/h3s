package loadbalancers

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/logger"
)

var UsePrivateIP = true

func Create(ctx clustercontext.ClusterContext, network *hcloud.Network, nodes []*hcloud.Server) *hcloud.LoadBalancer {
	balancer := Get(ctx)
	if balancer == nil {
		return create(ctx, network, nodes)
	}
	return balancer
}

func getNodeTarget(server *hcloud.Server) hcloud.LoadBalancerCreateOptsTarget {
	return hcloud.LoadBalancerCreateOptsTarget{
		Type:         hcloud.LoadBalancerTargetTypeServer,
		Server:       hcloud.LoadBalancerCreateOptsTargetServer{Server: server},
		UsePrivateIP: &UsePrivateIP,
	}
}

func create(
	ctx clustercontext.ClusterContext,
	network *hcloud.Network,
	nodes []*hcloud.Server,
) *hcloud.LoadBalancer {
	name := getName(ctx)

	logger.LogResourceEvent(logger.LoadBalancer, logger.Create, name, logger.Initialized)

	algorithm := hcloud.LoadBalancerAlgorithm{
		Type: "round_robin",
	}
	loadBalancerType := hcloud.LoadBalancerType{
		Name: "lb11",
	}
	location := hcloud.Location{Name: string(ctx.Config.ControlPlane.Pool.Location)}
	var targets []hcloud.LoadBalancerCreateOptsTarget
	for _, n := range nodes {
		targets = append(targets, getNodeTarget(n))
	}

	opts := hcloud.LoadBalancerCreateOpts{
		Name:             name,
		Targets:          targets,
		Location:         &location,
		Network:          network,
		Algorithm:        &algorithm,
		LoadBalancerType: &loadBalancerType,
		Labels:           ctx.GetLabels(),
	}

	res, _, err := ctx.Client.LoadBalancer.Create(ctx.Context, opts)

	if err != nil {
		logger.LogResourceEvent(logger.LoadBalancer, logger.Create, name, logger.Failure, err)
	}
	if res.LoadBalancer == nil {
		logger.LogResourceEvent(logger.LoadBalancer, logger.Create, name, logger.Failure, "Empty response")
	}
	if err := ctx.Client.Action.WaitFor(ctx.Context, res.Action); err != nil {
		logger.LogResourceEvent(logger.LoadBalancer, logger.Create, name, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.LoadBalancer, logger.Create, name, logger.Success)

	return res.LoadBalancer
}
