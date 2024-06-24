package install

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/k3s/install/command"
	"hcloud-k3s-cli/internal/resources/pool/node"
	"hcloud-k3s-cli/internal/resources/proxy"
	"hcloud-k3s-cli/internal/resources/server"
	"hcloud-k3s-cli/internal/utils/ssh"
)

func Install(
	ctx clustercontext.ClusterContext,
) {
	nodes := server.GetAll(ctx)

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

	for _, node := range controlPlaneNodes {
		cmd := command.ControlPlane(ctx, controlPlaneNodes, node)
		fmt.Printf("Installing controle plane k3s on %s\n", node.Name)
		ssh.Execute(ctx, p, node, cmd)
	}

	for _, node := range workerNodes {
		cmd := command.Worker(ctx, controlPlaneNodes)
		fmt.Printf("Installing worker k3s on %s\n", node.Name)
		ssh.Execute(ctx, p, node, cmd)
	}

	proxy.Delete(ctx)
}
