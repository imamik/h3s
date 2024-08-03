package loadbalancers

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/clustercontext"
	"h3s/internal/resources/network"
	"h3s/internal/utils/logger"
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

	balancer := Get(ctx)
	_, _, err = ctx.Client.RDNS.ChangeDNSPtr(ctx.Context, balancer, balancer.PublicNet.IPv4.IP, &ctx.Config.Domain)
	if err != nil {
		addEvent(logger.Failure, err)
		return nil
	}
	_, _, err = ctx.Client.RDNS.ChangeDNSPtr(ctx.Context, balancer, balancer.PublicNet.IPv6.IP, &ctx.Config.Domain)
	if err != nil {
		addEvent(logger.Failure, err)
		return nil
	}

	addEvent(logger.Success)
	return balancer
}
