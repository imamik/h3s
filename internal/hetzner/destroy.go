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
	"h3s/internal/utils/file"
	"h3s/internal/utils/logger"
	"sync"
)

func deleteResource(
	ctx *cluster.Cluster,
	wg *sync.WaitGroup,
	deleteFunc func(*cluster.Cluster),
) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		deleteFunc(ctx)
	}()
}

func Destroy(ctx *cluster.Cluster) error {
	logger.LogResourceEvent(logger.Cluster, logger.Delete, ctx.Config.Name, logger.Initialized)

	var wg sync.WaitGroup

	deleteResource(ctx, &wg, gateway.Delete)
	deleteResource(ctx, &wg, loadbalancers.Delete)
	deleteResource(ctx, &wg, pool.Delete)
	deleteResource(ctx, &wg, microos.Delete)
	deleteResource(ctx, &wg, dns.Delete)

	wg.Wait()             // Wait for all deletions to complete
	wg = sync.WaitGroup{} // Reset the wait group

	deleteResource(ctx, &wg, network.Delete)
	deleteResource(ctx, &wg, sshkey.Delete)

	wg.Wait() // Wait for all deletions to complete

	err := file.New("./k3s.yaml").Delete()
	if err != nil {
		return err // Ignore error
	}

	logger.LogResourceEvent(logger.Cluster, logger.Delete, ctx.Config.Name, logger.Success)
	return nil
}
