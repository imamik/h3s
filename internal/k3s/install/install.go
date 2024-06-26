package install

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/k3s/install/command"
	"hcloud-k3s-cli/internal/resources/loadbalancers/loadbalancer"
	"hcloud-k3s-cli/internal/resources/pool/node"
	"hcloud-k3s-cli/internal/resources/proxy"
	"hcloud-k3s-cli/internal/resources/server"
	"hcloud-k3s-cli/internal/utils/file"
	"hcloud-k3s-cli/internal/utils/logger"
	"hcloud-k3s-cli/internal/utils/ssh"
)

func downloadKubeConfig(
	ctx clustercontext.ClusterContext,
	proxy *hcloud.Server,
	remote *hcloud.Server,
) {
	cmd := "sudo cat /etc/rancher/k3s/k3s.yaml"
	kubeConfig, err := ssh.Execute(ctx, proxy, remote, cmd)
	if err != nil {
		logger.LogResourceEvent(logger.Server, "Download kubeconfig", remote.Name, logger.Failure, err)
	} else {
		err := file.Save([]byte(kubeConfig), "k3s.yaml")
		if err != nil {
			logger.LogResourceEvent(logger.Server, "Download kubeconfig", remote.Name, logger.Failure, err)
		} else {
			logger.LogResourceEvent(logger.Server, "Download kubeconfig", remote.Name, logger.Success)
		}
	}
}

func Install(ctx clustercontext.ClusterContext) {
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

	for i, remote := range controlPlaneNodes {
		cmd := command.ControlPlane(ctx, lb, controlPlaneNodes, remote)
		logger.LogResourceEvent(logger.Server, "Install Control Plane", remote.Name, logger.Initialized)
		ssh.Execute(ctx, proxyServer, remote, cmd)
		if i == 0 {
			downloadKubeConfig(ctx, proxyServer, remote)
		}
	}

	for _, remote := range workerNodes {
		cmd := command.Worker(ctx, lb)
		logger.LogResourceEvent(logger.Server, "Install Worker", remote.Name, logger.Initialized)
		ssh.Execute(ctx, proxyServer, remote, cmd)
	}

	proxy.Delete(ctx)
}
