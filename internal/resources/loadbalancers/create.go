package loadbalancers

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/network"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Create(ctx clustercontext.ClusterContext) *hcloud.LoadBalancer {
	balancer := Get(ctx)
	if balancer == nil {
		net := network.Get(ctx)
		return create(ctx, net)
	}
	return balancer
}

func create(
	ctx clustercontext.ClusterContext,
	net *hcloud.Network,
) *hcloud.LoadBalancer {
	name := getName(ctx)
	addEvent, logEvents := logger.NewEventLogger(logger.LoadBalancer, logger.Create, name)
	defer logEvents()

	algorithm := hcloud.LoadBalancerAlgorithm{
		Type: "round_robin",
	}
	loadBalancerType := hcloud.LoadBalancerType{
		Name: "lb11",
	}
	location := hcloud.Location{Name: string(ctx.Config.ControlPlane.Pool.Location)}

	opts := hcloud.LoadBalancerCreateOpts{
		Name:             name,
		Location:         &location,
		Network:          net,
		Algorithm:        &algorithm,
		LoadBalancerType: &loadBalancerType,
		Labels:           ctx.GetLabels(),
	}

	res, _, err := ctx.Client.LoadBalancer.Create(ctx.Context, opts)

	if err != nil {
		addEvent(logger.Failure, err)
		return nil
	}
	if res.LoadBalancer == nil {
		addEvent(logger.Failure, "Empty response")
		return nil
	}
	if err := ctx.Client.Action.WaitFor(ctx.Context, res.Action); err != nil {
		addEvent(logger.Failure, err)
		return nil
	}

	addEvent(logger.Success)
	return Get(ctx)
}
