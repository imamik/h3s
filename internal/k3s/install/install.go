package install

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/k3s/install/command"
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
	for _, node := range nodes {
		if node.Labels["is_control_plane"] == "true" {
			controlPlaneNodes = append(controlPlaneNodes, node)
		} else if node.Labels["is_worker"] == "true" {
			workerNodes = append(workerNodes, node)
		}
	}

	p := proxy.Create(ctx)

	for _, node := range controlPlaneNodes {
		cmd := command.ControlPlane(ctx, controlPlaneNodes, node)
		fmt.Printf("Installing controle plane k3s on %s\n", node.Name)
		ssh.Execute(ctx, p, node, cmd)
		ctx.Client.Server.Reboot(ctx.Context, node) // make sure to reboot the server after installation
	}

	for _, node := range workerNodes {
		cmd := command.Worker(ctx, controlPlaneNodes)
		fmt.Printf("Installing worker k3s on %s\n", node.Name)
		ssh.Execute(ctx, p, node, cmd)
		ctx.Client.Server.Reboot(ctx.Context, node) // make sure to reboot the server after installation
	}

	proxy.Delete(ctx)
}
