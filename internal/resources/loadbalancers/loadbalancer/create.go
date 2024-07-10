package loadbalancer

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Create(
	ctx clustercontext.ClusterContext,
	network *hcloud.Network,
	targets []hcloud.LoadBalancerCreateOptsTarget,
	balancerType Type,
) *hcloud.LoadBalancer {
	balancer := Get(ctx, balancerType)
	if balancer == nil {
		balancer = create(ctx, network, targets, balancerType)
	}
	return balancer
}

func create(
	ctx clustercontext.ClusterContext,
	network *hcloud.Network,
	targets []hcloud.LoadBalancerCreateOptsTarget,
	balancerType Type,
) *hcloud.LoadBalancer {
	name := getName(ctx, balancerType)

	logger.LogResourceEvent(logger.LoadBalancer, logger.Create, name, logger.Initialized)

	algorithm := hcloud.LoadBalancerAlgorithm{
		Type: "round_robin",
	}
	loadBalancerType := hcloud.LoadBalancerType{
		Name: "lb11",
	}
	labels := ctx.GetLabels(map[string]string{
		"type": string(balancerType),
	})
	location := hcloud.Location{Name: string(ctx.Config.ControlPlane.Pool.Location)}

	opts := hcloud.LoadBalancerCreateOpts{
		Name:             name,
		Targets:          targets,
		Location:         &location,
		Network:          network,
		Algorithm:        &algorithm,
		LoadBalancerType: &loadBalancerType,
		Labels:           labels,
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
