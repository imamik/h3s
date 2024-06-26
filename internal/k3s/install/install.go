package install

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/k3s/install/command"
	"hcloud-k3s-cli/internal/resources/loadbalancers/loadbalancer"
	"hcloud-k3s-cli/internal/resources/pool/node"
	"hcloud-k3s-cli/internal/resources/proxy"
	"hcloud-k3s-cli/internal/resources/server"
	"hcloud-k3s-cli/internal/utils/logger"
	"hcloud-k3s-cli/internal/utils/ssh"
)

func Install(
	ctx clustercontext.ClusterContext,
) {
	nodes := server.GetAll(ctx)

	balancerType := loadbalancer.ControlPlane
	if ctx.Config.CombinedLoadBalancer {
		balancerType = loadbalancer.Combined
	}
	lb := loadbalancer.Get(ctx, balancerType)

	var controlPlaneNodes []*hcloud.Server
	var workerNodes []*hcloud.Server
	for _, n := range nodes {
		if node.IsControlPlaneNode(n) {
			controlPlaneNodes = append(controlPlaneNodes, n)
		} else if node.IsWorkerNode(n) {
			workerNodes = append(workerNodes, n)
		}
	}

	p := proxy.Create(ctx)

	for _, n := range controlPlaneNodes {
		cmd := command.ControlPlane(ctx, lb, controlPlaneNodes, n)
		logger.LogResourceEvent(logger.Server, "Install Control Plane", n.Name, logger.Initialized)
		ssh.Execute(ctx, p, n, cmd)
	}

	for _, n := range workerNodes {
		cmd := command.Worker(ctx, lb)
		logger.LogResourceEvent(logger.Server, "Install Worker", n.Name, logger.Initialized)
		ssh.Execute(ctx, p, n, cmd)
	}

	proxy.Delete(ctx)
}
