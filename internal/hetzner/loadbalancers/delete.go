package loadbalancers

import (
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
)

func Delete(ctx *cluster.Cluster) {
	balancer := Get(ctx)
	if balancer == nil {
		return
	}

	addEvent, logEvents := logger.NewEventLogger(logger.LoadBalancer, logger.Delete, balancer.Name)
	defer logEvents()

	_, err := ctx.CloudClient.LoadBalancer.Delete(ctx.Context, balancer)
	if err != nil {
		addEvent(logger.Failure, err)
		return
	}

	addEvent(logger.Success)
}
