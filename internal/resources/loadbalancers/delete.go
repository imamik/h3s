package loadbalancers

import (
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Delete(ctx clustercontext.ClusterContext) {
	balancer := Get(ctx)
	if balancer == nil {
		return
	}

	logger.LogResourceEvent(logger.LoadBalancer, logger.Delete, balancer.Name, logger.Initialized)

	_, err := ctx.Client.LoadBalancer.Delete(ctx.Context, balancer)
	if err != nil {
		logger.LogResourceEvent(logger.LoadBalancer, logger.Delete, balancer.Name, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.LoadBalancer, logger.Delete, balancer.Name, logger.Success)
}
