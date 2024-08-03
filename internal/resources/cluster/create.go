package cluster

import (
	"h3s/internal/clustercontext"
	"h3s/internal/resources/dns"
	"h3s/internal/resources/gateway"
	"h3s/internal/resources/loadbalancers"
	"h3s/internal/resources/microos"
	"h3s/internal/resources/network"
	"h3s/internal/resources/pool"
	"h3s/internal/resources/sshkey"
	"h3s/internal/utils/logger"
	"sync"
)

func Create(ctx clustercontext.ClusterContext) {
	addEvent, logEvents := logger.NewEventLogger(logger.Cluster, logger.Create, ctx.Config.Name)
	addEvent(logger.Initialized)
	defer logEvents()

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		microos.Create(ctx)
	}()
	go func() {
		sshkey.Create(ctx)
		wg.Done()
	}()
	go func() {
		network.Create(ctx)
		wg.Done()
	}()
	wg.Wait()

	wg = sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		loadbalancers.Create(ctx)
		dns.Create(ctx)
	}()
	go func() {
		defer wg.Done()
		pool.CreatePools(ctx)
	}()
	if ctx.Config.PublicIps == false {
		wg.Add(1)
		go func() {
			defer wg.Done()
			gateway.Create(ctx)
		}()
	}
	wg.Wait()

	addEvent(logger.Success)
}
