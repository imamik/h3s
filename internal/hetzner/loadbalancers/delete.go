package loadbalancers

import (
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
)

// Delete deletes the Hetzner cloud load balancer
func Delete(ctx *cluster.Cluster) error {
	l := logger.New(nil, logger.LoadBalancer, logger.Delete, "")
	defer l.LogEvents()

	balancer, err := Get(ctx)
	if balancer == nil && err.Error() == "load balancer is nil" {
		l.AddEvent(logger.Success)
		return nil
	}

	_, err = ctx.CloudClient.LoadBalancer.Delete(ctx.Context, balancer)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	l.AddEvent(logger.Success)
	return nil
}
