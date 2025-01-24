// Package loadbalancers contains the functionality for creating a Hetzner cloud load balancer
package loadbalancers

import (
	"errors"
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Create sets up a new load balancer in the Hetzner Cloud for the cluster
func Create(ctx *cluster.Cluster, net *hcloud.Network) (*hcloud.LoadBalancer, error) {
	balancer, err := Get(ctx)
	if balancer != nil && err == nil {
		return balancer, nil
	}
	return create(ctx, net)
}

func create(
	ctx *cluster.Cluster,
	net *hcloud.Network,
) (*hcloud.LoadBalancer, error) {
	name := getName(ctx)

	l := logger.New(nil, logger.LoadBalancer, logger.Create, name)
	defer l.LogEvents()

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
		l.AddEvent(logger.Failure, err)
		return nil, err
	}
	if res.LoadBalancer == nil {
		err = errors.New("load balancer is nil")
		l.AddEvent(logger.Failure, err)
		return nil, err
	}
	if waitErr := ctx.CloudClient.Action.WaitFor(ctx.Context, res.Action); waitErr != nil {
		l.AddEvent(logger.Failure, waitErr)
		return nil, waitErr
	}

	balancer, err := Get(ctx)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	_, _, err = ctx.CloudClient.RDNS.ChangeDNSPtr(ctx.Context, balancer, balancer.PublicNet.IPv4.IP, &ctx.Config.Domain)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}
	_, _, err = ctx.CloudClient.RDNS.ChangeDNSPtr(ctx.Context, balancer, balancer.PublicNet.IPv6.IP, &ctx.Config.Domain)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	l.AddEvent(logger.Success)
	return balancer, nil
}
