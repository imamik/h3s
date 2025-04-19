// Package pool contains the functionality for managing Hetzner cloud pools
package pool

import (
	"fmt"
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/hetzner/microos"
	"h3s/internal/hetzner/pool/node"
	"h3s/internal/hetzner/pool/placementgroup"
	"h3s/internal/utils/logger"
	"sync"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// CreatePools creates the control plane and worker pools
func CreatePools(
	ctx *cluster.Cluster,
	sshKey *hcloud.SSHKey,
	net *hcloud.Network,
	images *microos.ImageInArchitecture,
) ([]*hcloud.Server, error) {
	// Defensive: nil checks for required arguments
	if ctx == nil {
		return nil, fmt.Errorf("CreatePools: ctx is nil")
	}
	if sshKey == nil {
		return nil, fmt.Errorf("CreatePools: sshKey is nil")
	}
	if net == nil {
		return nil, fmt.Errorf("CreatePools: net is nil")
	}
	if images == nil {
		return nil, fmt.Errorf("CreatePools: images is nil")
	}

	l := logger.New(nil, logger.Pool, logger.Create, "All")
	defer l.LogEvents()

	// Create a channel to collect the nodes & setup a WaitGroup
	nodeCh := make(chan []*hcloud.Server)
	errCh := make(chan error)
	var wg sync.WaitGroup

	// Create control plane pool in a goroutine
	wg.Add(1)
	go func() {
		logr := logger.New(l, logger.Pool, logger.Create, "Control Plane")
		defer wg.Done()
		nodes, err := CreatePool(
			ctx,
			sshKey,
			net,
			ctx.Config.ControlPlane.Pool,
			true,
			ctx.Config.ControlPlane.AsWorkerPool,
			images,
		)
		if err != nil {
			logr.AddEvent(logger.Failure, err)
			errCh <- err
			return
		}
		nodeCh <- nodes
	}()

	// Create worker pools in separate goroutines
	for _, pool := range ctx.Config.WorkerPools {
		logr := logger.New(l, logger.Pool, logger.Create, ctx.GetName(pool.Name))
		wg.Add(1)
		go func(pool config.NodePool) {
			defer wg.Done()
			nodes, err := CreatePool(
				ctx,
				sshKey,
				net,
				pool,
				false,
				true,
				images,
			)
			if err != nil {
				logr.AddEvent(logger.Failure, err)
				errCh <- err
				return
			}
			nodeCh <- nodes
		}(pool)
	}

	// Wait for all goroutines to finish & close the channel
	go func() {
		wg.Wait()
		close(nodeCh)
		close(errCh)
	}()

	// Collect the nodes from the channel
	var nodes []*hcloud.Server
	for n := range nodeCh {
		nodes = append(nodes, n...)
	}

	// Check and handle errors
	//nolint:prealloc // Suppressing pre-allocation suggestion intentionally as we don't know the number of errors
	var errors []error
	for err := range errCh {
		errors = append(errors, err)
	}
	if len(errors) > 0 {
		l.AddEvent(logger.Failure, errors)
		return nil, fmt.Errorf("failed to create pools: %v", errors)
	}

	l.AddEvent(logger.Success)
	return nodes, nil
}

// CreatePool creates a pool of Hetzner cloud servers
func CreatePool(
	ctx *cluster.Cluster,
	sshKey *hcloud.SSHKey,
	network *hcloud.Network,
	pool config.NodePool,
	isControlPlane bool,
	isWorker bool,
	images *microos.ImageInArchitecture,
) ([]*hcloud.Server, error) {
	l := logger.New(nil, logger.Pool, logger.Create, ctx.GetName(pool.Name))
	defer l.LogEvents()

	// Defensive: nil checks for ctx, images, and pool.Instance
	if ctx == nil {
		return nil, fmt.Errorf("CreatePool: ctx is nil")
	}
	if images == nil {
		return nil, fmt.Errorf("CreatePool: images is nil")
	}
	if pool.Instance == "" {
		return nil, fmt.Errorf("CreatePool: pool.Instance is empty")
	}

	img := images.ARM
	if config.GetArchitecture(pool.Instance) == hcloud.ArchitectureX86 {
		img = images.X86
	}

	placementGroup, err := placementgroup.Create(ctx, pool, isControlPlane, isWorker)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	// Create a channel to collect the nodes & setup a WaitGroup
	nodeCh := make(chan *hcloud.Server, pool.Nodes)
	errCh := make(chan error)
	var wg sync.WaitGroup

	for i := 0; i < pool.Nodes; i++ {
		wg.Add(1) // Increment the WaitGroup counter
		go func(i int) {
			defer wg.Done() // Decrement the counter when the goroutine completes
			n, err := node.Create(
				ctx,
				sshKey,
				network,
				img,
				placementGroup,
				pool,
				i,
				isControlPlane,
				isWorker,
			)
			if err != nil {
				l.AddEvent(logger.Failure, err)
				errCh <- err
				return
			}
			nodeCh <- n
		}(i)
	}

	// Wait for all goroutines to finish & close the channel
	wg.Wait()
	close(nodeCh)
	close(errCh)

	// Check for errors
	if err, ok := <-errCh; ok {
		return nil, err
	}

	// Collect the nodes from the channel
	nodes := make([]*hcloud.Server, 0, pool.Nodes)
	for n := range nodeCh {
		nodes = append(nodes, n)
	}

	l.AddEvent(logger.Success)
	return nodes, nil
}
