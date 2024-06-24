package loadbalancer

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/clustercontext"
	"hcloud-k3s-cli/pkg/utils/logger"
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

	balancer, _, err := ctx.Client.LoadBalancer.Create(ctx.Context, hcloud.LoadBalancerCreateOpts{
		Name:        name,
		NetworkZone: ctx.Config.NetworkZone,
		Network:     network,
		Algorithm: &hcloud.LoadBalancerAlgorithm{
			Type: "round_robin",
		},
		LoadBalancerType: &hcloud.LoadBalancerType{
			Name: "lb11",
		},
		Labels: ctx.GetLabels(map[string]string{
			"type": string(balancerType),
		}),
	})
	if err != nil || balancer.LoadBalancer == nil {
		logger.LogResourceEvent(logger.LoadBalancer, logger.Create, name, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.LoadBalancer, logger.Create, name, logger.Success)

	return balancer.LoadBalancer
}
