package pool

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/config"
	"hcloud-k3s-cli/internal/resources/pool/node"
	"hcloud-k3s-cli/internal/resources/pool/placementgroup"
	"hcloud-k3s-cli/internal/utils/logger"
)

func CreatePools(
	ctx clustercontext.ClusterContext,
	sshKey *hcloud.SSHKey,
	network *hcloud.Network,
) []*hcloud.Server {
	// Create control plane pool
	nodes := create(
		ctx,
		sshKey,
		network,
		ctx.Config.ControlPlane.Pool,
		true,
		ctx.Config.ControlPlane.AsWorkerPool,
	)

	// Create worker pools
	for _, pool := range ctx.Config.WorkerPools {
		workerNodes := create(
			ctx,
			sshKey,
			network,
			pool,
			false,
			true,
		)
		nodes = append(nodes, workerNodes...)
	}

	return nodes
}

func create(
	ctx clustercontext.ClusterContext,
	sshKey *hcloud.SSHKey,
	network *hcloud.Network,
	pool config.NodePool,
	isControlPlane bool,
	isWorker bool,
) []*hcloud.Server {
	logger.LogResourceEvent(logger.Pool, logger.Create, ctx.GetName(pool.Name), logger.Initialized)

	placementGroup := placementgroup.Create(ctx, pool, isControlPlane, isWorker)
	var nodes []*hcloud.Server

	for i := 0; i < pool.Nodes; i++ {
		n := node.Create(
			ctx,
			sshKey,
			network,
			placementGroup,
			pool,
			i,
			isControlPlane,
			isWorker,
		)
		nodes = append(nodes, n)
	}

	logger.LogResourceEvent(logger.Pool, logger.Create, ctx.GetName(pool.Name), logger.Success)
	return nodes
}
