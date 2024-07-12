package loadbalancers

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Get(ctx clustercontext.ClusterContext) *hcloud.LoadBalancer {
	balancer := getName(ctx)
	logger.LogResourceEvent(logger.LoadBalancer, logger.Get, balancer, logger.Initialized)

	lb, _, err := ctx.Client.LoadBalancer.GetByName(ctx.Context, balancer)
	if err != nil || lb == nil {
		logger.LogResourceEvent(logger.LoadBalancer, logger.Get, balancer, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.LoadBalancer, logger.Get, balancer, logger.Success)
	return lb
}
