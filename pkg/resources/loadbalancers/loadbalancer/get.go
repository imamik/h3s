package loadbalancer

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/clustercontext"
	"hcloud-k3s-cli/pkg/utils/logger"
)

func Get(ctx clustercontext.ClusterContext, balancerType Type) *hcloud.LoadBalancer {
	balancer := getName(ctx, balancerType)
	logger.LogResourceEvent(logger.LoadBalancer, logger.Get, balancer, logger.Initialized)

	network, _, err := ctx.Client.LoadBalancer.GetByName(ctx.Context, balancer)
	if err != nil || network == nil {
		logger.LogResourceEvent(logger.LoadBalancer, logger.Get, balancer, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.LoadBalancer, logger.Get, balancer, logger.Success)
	return network
}
