package loadbalancers

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/hetzner/network"
	"h3s/internal/utils/logger"
)

func Create(ctx *cluster.Cluster) *hcloud.LoadBalancer {
	balancer := Get(ctx)
	if balancer == nil {
		net := network.Get(ctx)
		return create(ctx, net)
	}
	return balancer
}

func create(
	ctx *cluster.Cluster,
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

	res, _, err := ctx.CloudClient.LoadBalancer.Create(ctx.Context, opts)

	if err != nil {
		addEvent(logger.Failure, err)
		return nil
	}
	if res.LoadBalancer == nil {
		addEvent(logger.Failure, "Empty response")
		return nil
	}
	if err := ctx.CloudClient.Action.WaitFor(ctx.Context, res.Action); err != nil {
		addEvent(logger.Failure, err)
		return nil
	}

	balancer := Get(ctx)
	_, _, err = ctx.CloudClient.RDNS.ChangeDNSPtr(ctx.Context, balancer, balancer.PublicNet.IPv4.IP, &ctx.Config.Domain)
	if err != nil {
		addEvent(logger.Failure, err)
		return nil
	}
	_, _, err = ctx.CloudClient.RDNS.ChangeDNSPtr(ctx.Context, balancer, balancer.PublicNet.IPv6.IP, &ctx.Config.Domain)
	if err != nil {
		addEvent(logger.Failure, err)
		return nil
	}

	addEvent(logger.Success)
	return balancer
}
