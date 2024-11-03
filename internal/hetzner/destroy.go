package hetzner

import (
	"h3s/internal/cluster"
	"h3s/internal/hetzner/dns"
	"h3s/internal/hetzner/gateway"
	"h3s/internal/hetzner/loadbalancers"
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
	deleteFunc func(*cluster.Cluster) error,
) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := deleteFunc(ctx)
		if err != nil {
			return
		}
	}()
}

// Destroy deletes the Hetzner cloud cluster
func Destroy(ctx *cluster.Cluster) error {
	l := logger.New(nil, logger.Cluster, logger.Delete, ctx.Config.Name)
	defer l.LogEvents()

	var wg sync.WaitGroup

	deleteResource(ctx, &wg, gateway.Delete)
	deleteResource(ctx, &wg, loadbalancers.Delete)
	deleteResource(ctx, &wg, pool.Delete)
	// deleteResource(ctx, &wg, microos.Delete)
	deleteResource(ctx, &wg, dns.Delete)

	wg.Wait()             // Wait for all deletions to complete
	wg = sync.WaitGroup{} // Reset the wait group

	deleteResource(ctx, &wg, network.Delete)
	deleteResource(ctx, &wg, sshkey.Delete)

	wg.Wait() // Wait for all deletions to complete

	err := file.New("./h3s-kubeconfig.yaml").Delete()
	if err != nil {
		return err // Ignore error
	}

	l.AddEvent(logger.Success)
	return nil
}
