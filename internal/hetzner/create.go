// Package hetzner contains the functionality for creating the Hetzner cloud cluster
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
)

// Create creates the Hetzner cloud cluster
func Create(ctx *cluster.Cluster) error {
	l := logger.New(nil, logger.Cluster, logger.Create, ctx.Config.Name)
	defer l.LogEvents()

	sshKey, err := sshkey.Create(ctx)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	images, err := microos.Create(ctx, sshKey)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	net, err := network.Create(ctx)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	lb, err := loadbalancers.Create(ctx, net)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	if _, err := dns.Create(ctx, lb); err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	if _, err := pool.CreatePools(ctx, sshKey, net, images); err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	if _, err := gateway.Create(ctx, sshKey, net, images); err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	l.AddEvent(logger.Success)
	return nil
}
