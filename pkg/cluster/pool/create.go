package pool

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/cluster/pool/placementgroup"
	"hcloud-k3s-cli/pkg/cluster/pool/server"
	"hcloud-k3s-cli/pkg/config"
	"hcloud-k3s-cli/pkg/utils/logger"
)

func CreatePools(
	ctx clustercontext.ClusterContext,
	sshKey *hcloud.SSHKey,
	network *hcloud.Network,
) {
	// Create control plane pool
	create(
		ctx,
		sshKey,
		network,
		ctx.Config.ControlPlane.Pool,
		true,
		ctx.Config.ControlPlane.AsWorkerPool,
	)

	// Create worker pools
	for _, pool := range ctx.Config.WorkerPools {
		create(
			ctx,
			sshKey,
			network,
			pool,
			false,
			true,
		)
	}
}

func create(
	ctx clustercontext.ClusterContext,
	sshKey *hcloud.SSHKey,
	network *hcloud.Network,
	pool config.NodePool,
	isControlPlane bool,
	isWorker bool,
) {
	logger.LogResourceEvent(logger.Pool, logger.Create, ctx.GetName(pool.Name), logger.Initialized)

	placementGroup := placementgroup.Create(ctx, pool)

	for i := 0; i < pool.Nodes; i++ {
		server.Create(
			ctx,
			sshKey,
			network,
			placementGroup,
			pool,
			i,
			isControlPlane,
			isWorker,
		)
	}

	logger.LogResourceEvent(logger.Pool, logger.Create, ctx.GetName(pool.Name), logger.Success)
}
