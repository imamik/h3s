package cluster

import (
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/dns"
	"hcloud-k3s-cli/internal/resources/gateway"
	"hcloud-k3s-cli/internal/resources/loadbalancers"
	"hcloud-k3s-cli/internal/resources/microos"
	"hcloud-k3s-cli/internal/resources/network"
	"hcloud-k3s-cli/internal/resources/pool"
	"hcloud-k3s-cli/internal/resources/sshkey"
	"hcloud-k3s-cli/internal/utils/logger"
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

	loadbalancers.Add(ctx)

	addEvent(logger.Success)
}
