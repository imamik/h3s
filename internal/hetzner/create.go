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
		if err := microos.Create(ctx); err != nil {
			l.AddEvent(logger.Failure, err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, err := sshkey.Create(ctx); err != nil {
			l.AddEvent(logger.Failure, err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, err := network.Create(ctx); err != nil {
			l.AddEvent(logger.Failure, err)
		}
	}()
	wg.Wait()

	wg = sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, err := loadbalancers.Create(ctx); err != nil {
			l.AddEvent(logger.Failure, err)
		}
		if err := dns.Create(ctx); err != nil {
			l.AddEvent(logger.Failure, err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, err := pool.CreatePools(ctx); err != nil {
			l.AddEvent(logger.Failure, err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, err := gateway.Create(ctx); err != nil {
			l.AddEvent(logger.Failure, err)
		}
	}()

	l.AddEvent(logger.Success)
}
