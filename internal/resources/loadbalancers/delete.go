package loadbalancers

import (
	"h3s/internal/clustercontext"
	"h3s/internal/utils/logger"
)

func Delete(ctx clustercontext.ClusterContext) {
	balancer := Get(ctx)
	if balancer == nil {
		return
	}

	addEvent, logEvents := logger.NewEventLogger(logger.LoadBalancer, logger.Delete, balancer.Name)
	defer logEvents()

	_, err := ctx.Client.LoadBalancer.Delete(ctx.Context, balancer)
	if err != nil {
		addEvent(logger.Failure, err)
		return
	}

	addEvent(logger.Success)
}
