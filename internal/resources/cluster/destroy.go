package cluster

import (
	"h3s/internal/clustercontext"
	"h3s/internal/resources/dns"
	"h3s/internal/resources/gateway"
	"h3s/internal/resources/loadbalancers"
	"h3s/internal/resources/network"
	"h3s/internal/resources/pool"
	"h3s/internal/resources/sshkey"
	"h3s/internal/utils/file"
	"h3s/internal/utils/logger"
	"sync"
)

func deleteResource(
	ctx clustercontext.ClusterContext,
	wg *sync.WaitGroup,
	deleteFunc func(clustercontext.ClusterContext),
) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		deleteFunc(ctx)
	}()
}

func Destroy(ctx clustercontext.ClusterContext) {
	logger.LogResourceEvent(logger.Cluster, logger.Delete, ctx.Config.Name, logger.Initialized)

	var wg sync.WaitGroup

	deleteResource(ctx, &wg, gateway.Delete)
	deleteResource(ctx, &wg, loadbalancers.Delete)
	deleteResource(ctx, &wg, pool.Delete)
	deleteResource(ctx, &wg, dns.Delete)

	wg.Wait()             // Wait for all deletions to complete
	wg = sync.WaitGroup{} // Reset the wait group

	deleteResource(ctx, &wg, network.Delete)
	deleteResource(ctx, &wg, sshkey.Delete)

	wg.Wait() // Wait for all deletions to complete

	file.Delete("./k3s.yaml")

	logger.LogResourceEvent(logger.Cluster, logger.Delete, ctx.Config.Name, logger.Success)
}
