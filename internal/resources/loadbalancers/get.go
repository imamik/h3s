package loadbalancers

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/clustercontext"
	"h3s/internal/utils/logger"
)

func Get(ctx clustercontext.ClusterContext) *hcloud.LoadBalancer {
	balancer := getName(ctx)

	addEvent, logEvents := logger.NewEventLogger(logger.LoadBalancer, logger.Create, balancer)
	defer logEvents()

	lb, _, err := ctx.Client.LoadBalancer.GetByName(ctx.Context, balancer)
	if err != nil || lb == nil {
		addEvent(logger.Failure, err)
		return nil
	}

	addEvent(logger.Success)
	return lb
}
