package install

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/k3s/install/commands"
	"hcloud-k3s-cli/internal/k3s/install/software"
	"hcloud-k3s-cli/internal/resources/loadbalancers/loadbalancer"
	"hcloud-k3s-cli/internal/resources/network"
	"hcloud-k3s-cli/internal/resources/pool/node"
	"hcloud-k3s-cli/internal/resources/proxy"
	"hcloud-k3s-cli/internal/resources/server"
	"hcloud-k3s-cli/internal/utils/logger"
	"hcloud-k3s-cli/internal/utils/ssh"
)

func getSetup(ctx clustercontext.ClusterContext) (*hcloud.Network, *hcloud.LoadBalancer, *hcloud.Server, []*hcloud.Server, []*hcloud.Server) {
	net := network.Get(ctx)
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

	proxyServer := proxy.Create(ctx)

	return net, lb, proxyServer, controlPlaneNodes, workerNodes
}

func Install(ctx clustercontext.ClusterContext, cleanup bool) {
	net, lb, proxyServer, controlPlaneNodes, workerNodes := getSetup(ctx)

	if cleanup {
		defer proxy.Delete(ctx)
	}

	for i, remote := range controlPlaneNodes {
		cmd := commands.ControlPlane(ctx, lb, controlPlaneNodes, remote)
		logger.LogResourceEvent(logger.Server, "Install Control Plane", remote.Name, logger.Initialized)
		ssh.ExecuteViaProxy(ctx, proxyServer, remote, cmd)
		if i == 0 {
			downloadKubeConfig(ctx, lb, proxyServer, remote)
			software.Install(ctx, net, lb, proxyServer, remote)
		}
	}

	for _, remote := range workerNodes {
		cmd := commands.Worker(ctx, lb, controlPlaneNodes, remote)
		logger.LogResourceEvent(logger.Server, "Install Worker", remote.Name, logger.Initialized)
		ssh.ExecuteViaProxy(ctx, proxyServer, remote, cmd)
	}

}

func InstallSoftware(ctx clustercontext.ClusterContext, cleanup bool) {
	net, lb, proxyServer, controlPlaneNodes, _ := getSetup(ctx)

	if cleanup {
		defer proxy.Delete(ctx)
	}

	software.Install(ctx, net, lb, proxyServer, controlPlaneNodes[0])
}
