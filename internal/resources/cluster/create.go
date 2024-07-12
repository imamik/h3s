package cluster

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/dns"
	"hcloud-k3s-cli/internal/resources/gateway"
	"hcloud-k3s-cli/internal/resources/loadbalancers"
	"hcloud-k3s-cli/internal/resources/microos"
	"hcloud-k3s-cli/internal/resources/network"
	"hcloud-k3s-cli/internal/resources/pool"
	"hcloud-k3s-cli/internal/resources/server"
	"hcloud-k3s-cli/internal/resources/sshkey"
	"hcloud-k3s-cli/internal/utils/logger"
	"sync"
)

func Create(ctx clustercontext.ClusterContext) {
	logger.LogResourceEvent(logger.Cluster, logger.Create, ctx.Config.Name, logger.Initialized)

	var sshKey *hcloud.SSHKey
	var net *hcloud.Network
	var lb *hcloud.LoadBalancer

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		microos.Create(ctx)
	}()
	go func() {
		sshKey = sshkey.Create(ctx)
		wg.Done()
	}()
	go func() {
		net = network.Create(ctx)
		wg.Done()
	}()
	wg.Wait()

	wg = sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		lb = loadbalancers.Create(ctx, net)
		dns.Create(ctx, lb)
	}()
	go func() {
		defer wg.Done()
		pool.CreatePools(ctx, sshKey, net)
	}()
	if ctx.Config.PublicIps == false {
		wg.Add(1)
		go func() {
			defer wg.Done()
			gateway.Create(ctx)
		}()
	}
	wg.Wait()

	nodes := server.GetAll(ctx)
	loadbalancers.Add(ctx, lb, nodes)

	logger.LogResourceEvent(logger.Cluster, logger.Create, ctx.Config.Name, logger.Success)
}
