package loadbalancers

import (
	"errors"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
)

func Get(ctx *cluster.Cluster) (*hcloud.LoadBalancer, error) {
	balancer := getName(ctx)

	l := logger.New(nil, logger.LoadBalancer, logger.Get, balancer)
	defer l.LogEvents()

	lb, _, err := ctx.CloudClient.LoadBalancer.GetByName(ctx.Context, balancer)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}
	if lb == nil {
		err := errors.New("load balancer is nil")
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	l.AddEvent(logger.Success)
	return lb, nil
}