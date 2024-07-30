package install

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/k3s/install/commands"
	"hcloud-k3s-cli/internal/k3s/install/software"
	"hcloud-k3s-cli/internal/k3s/kubeconfig"
	"hcloud-k3s-cli/internal/resources/gateway"
	"hcloud-k3s-cli/internal/resources/loadbalancers"
	"hcloud-k3s-cli/internal/resources/network"
	"hcloud-k3s-cli/internal/resources/pool/node"
	"hcloud-k3s-cli/internal/resources/server"
	"sort"
)

func getSetup(ctx clustercontext.ClusterContext) (*hcloud.Network, *hcloud.LoadBalancer, *hcloud.Server, []*hcloud.Server, []*hcloud.Server) {
	net := network.Get(ctx)
	nodes := server.GetAll(ctx)
	lb := loadbalancers.Get(ctx)

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Name < nodes[j].Name
	})

	var controlPlaneNodes []*hcloud.Server
	var workerNodes []*hcloud.Server
	for _, n := range nodes {
		if node.IsControlPlaneNode(n) {
			controlPlaneNodes = append(controlPlaneNodes, n)
		} else if node.IsWorkerNode(n) {
			workerNodes = append(workerNodes, n)
		}
	}

	var gatewayServer *hcloud.Server
	if ctx.Config.PublicIps == false {
		gatewayServer, _ = gateway.Get(ctx)
	}

	return net, lb, gatewayServer, controlPlaneNodes, workerNodes
}

func Install(ctx clustercontext.ClusterContext) {
	net, lb, gatewayServer, controlPlaneNodes, workerNodes := getSetup(ctx)
	firstControlPlane := controlPlaneNodes[0]

	for _, remote := range controlPlaneNodes {
		commands.ControlPlane(ctx, lb, controlPlaneNodes, gatewayServer, remote)
	}

	for _, remote := range workerNodes {
		commands.Worker(ctx, controlPlaneNodes[0], gatewayServer, remote)
	}

	software.Install(ctx, net, lb, gatewayServer, firstControlPlane)
	kubeconfig.DownloadKubeConfig(ctx, gatewayServer, firstControlPlane)

}

func InstallSoftware(ctx clustercontext.ClusterContext) {
	net, lb, proxyServer, controlPlaneNodes, _ := getSetup(ctx)

	software.Install(ctx, net, lb, proxyServer, controlPlaneNodes[0])
}
