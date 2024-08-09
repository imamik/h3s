package install

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/clustercontext"
	"h3s/internal/k3s/install/commands"
	"h3s/internal/k3s/install/software"
	"h3s/internal/k3s/kubeconfig"
	"h3s/internal/resources/gateway"
	"h3s/internal/resources/loadbalancers"
	"h3s/internal/resources/network"
	"h3s/internal/resources/pool/node"
	"h3s/internal/resources/server"
	"sort"
)

func GetSetup(ctx clustercontext.ClusterContext) (*hcloud.Network, *hcloud.LoadBalancer, *hcloud.Server, []*hcloud.Server, []*hcloud.Server) {
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

func K3s(ctx clustercontext.ClusterContext) {
	_, lb, gatewayServer, controlPlaneNodes, workerNodes := GetSetup(ctx)

	for _, remote := range controlPlaneNodes {
		commands.ControlPlane(ctx, lb, controlPlaneNodes, gatewayServer, remote)
	}

	for _, remote := range workerNodes {
		commands.Worker(ctx, controlPlaneNodes[0], gatewayServer, remote)
	}

}

func Software(ctx clustercontext.ClusterContext) {
	net, lb, proxyServer, controlPlaneNodes, _ := GetSetup(ctx)

	software.Install(ctx, net, lb, proxyServer, controlPlaneNodes[0])
}

func DownloadKubeconfig(ctx clustercontext.ClusterContext) {
	_, _, gatewayServer, controlPlaneNodes, _ := GetSetup(ctx)
	firstControlPlane := controlPlaneNodes[0]

	kubeconfig.Download(ctx, gatewayServer, firstControlPlane)
}
