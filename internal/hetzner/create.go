package hetzner

import (
	"h3s/internal/cluster"
	"h3s/internal/hetzner/dns"
	"h3s/internal/hetzner/gateway"
	"h3s/internal/hetzner/loadbalancers"
	"h3s/internal/hetzner/microos"
	"h3s/internal/hetzner/network"
	"h3s/internal/hetzner/pool"
	"h3s/internal/hetzner/sshkey"
	"h3s/internal/utils/logger"
	"sync"
)

func Create(ctx *cluster.Cluster) {
	l := logger.New(nil, logger.Cluster, logger.Create, ctx.Config.Name)
	defer l.LogEvents()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		microos.Create(ctx)
	}()
	wg.Add(1)
	go func() {
		sshkey.Create(ctx)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		network.Create(ctx)
		wg.Done()
	}()
	wg.Wait()

	wg = sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		loadbalancers.Create(ctx)
		dns.Create(ctx)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		pool.CreatePools(ctx)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		gateway.Create(ctx)
	}()

	l.AddEvent(logger.Success)
}