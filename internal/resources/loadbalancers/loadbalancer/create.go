package loadbalancer

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Create(ctx clustercontext.ClusterContext, network *hcloud.Network, balancerType Type) *hcloud.LoadBalancer {
	balancer := Get(ctx, balancerType)
	if balancer == nil {
		balancer = create(ctx, network, balancerType)
	}
	return balancer
}

func create(ctx clustercontext.ClusterContext, network *hcloud.Network, balancerType Type) *hcloud.LoadBalancer {
	name := getName(ctx, balancerType)

	logger.LogResourceEvent(logger.LoadBalancer, logger.Create, name, logger.Initialized)

	services := []hcloud.LoadBalancerCreateOptsService{
		{
			Protocol:        hcloud.LoadBalancerServiceProtocolHTTP,
			ListenPort:      hcloud.Ptr(80),
			DestinationPort: hcloud.Ptr(80),
		},
		{
			Protocol:        hcloud.LoadBalancerServiceProtocolTCP,
			ListenPort:      hcloud.Ptr(6443),
			DestinationPort: hcloud.Ptr(6443),
		},
	}
	algorithm := hcloud.LoadBalancerAlgorithm{
		Type: "round_robin",
	}
	loadBalancerType := hcloud.LoadBalancerType{
		Name: "lb11",
	}
	labels := ctx.GetLabels(map[string]string{
		"type": string(balancerType),
	})
	location := &hcloud.Location{Name: string(ctx.Config.ControlPlane.Pool.Location)}

	balancer, _, err := ctx.Client.LoadBalancer.Create(ctx.Context, hcloud.LoadBalancerCreateOpts{
		Name:             name,
		NetworkZone:      ctx.Config.NetworkZone,
		Network:          network,
		Location:         location,
		Algorithm:        &algorithm,
		LoadBalancerType: &loadBalancerType,
		Services:         services,
		Labels:           labels,
	})
	if err != nil || balancer.LoadBalancer == nil {
		logger.LogResourceEvent(logger.LoadBalancer, logger.Create, name, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.LoadBalancer, logger.Create, name, logger.Success)

	return balancer.LoadBalancer
}
