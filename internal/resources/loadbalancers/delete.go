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

	addEvent, logEvents := logger.NewEventLogger(logger.LoadBalancer, logger.Create, balancer.Name)
	defer logEvents()

	_, err := ctx.Client.LoadBalancer.Delete(ctx.Context, balancer)
	if err != nil {
		addEvent(logger.Failure, err)
		return
	}

	addEvent(logger.Success)
}
