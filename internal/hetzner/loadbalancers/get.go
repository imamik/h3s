package loadbalancers

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
)

func Get(ctx *cluster.Cluster) *hcloud.LoadBalancer {
	balancer := getName(ctx)

	addEvent, logEvents := logger.NewEventLogger(logger.LoadBalancer, logger.Create, balancer)
	defer logEvents()

	lb, _, err := ctx.CloudClient.LoadBalancer.GetByName(ctx.Context, balancer)
	if err != nil || lb == nil {
		addEvent(logger.Failure, err)
		return nil
	}

	addEvent(logger.Success)
	return lb
}
