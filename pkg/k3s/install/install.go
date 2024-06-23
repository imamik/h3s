package install

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/melbahja/goph"
	"hcloud-k3s-cli/pkg/clustercontext"
	"hcloud-k3s-cli/pkg/k3s/install/command"
	"hcloud-k3s-cli/pkg/resources/gateway"
	"hcloud-k3s-cli/pkg/resources/server"
	"hcloud-k3s-cli/pkg/utils/ssh"
)

func Install(
	ctx clustercontext.ClusterContext,
) {
	nodes := server.GetAllByProject(ctx)
	if len(nodes) == 0 {
		fmt.Println("No servers found")
		return
	}

	var controlPlaneNodes []*hcloud.Server
	var workerNodes []*hcloud.Server
	for _, node := range nodes {
		if node.Labels["is_control_plane"] == "true" {
			controlPlaneNodes = append(controlPlaneNodes, node)
		} else {
			workerNodes = append(workerNodes, node)
		}
	}

	gate := gateway.Create(ctx)
	client := ssh.Client(ctx, gate)

	defer func(client *goph.Client) {
		err := client.Close()
		if err != nil {
			return
		}
	}(client)

	for _, node := range controlPlaneNodes {
		cmd := command.ControlPlane(ctx, controlPlaneNodes, node)
		fmt.Printf("Installing controle plane k3s on %s\n", node.Name)
		ssh.Execute(ctx, client, node, cmd)
	}

	for _, node := range workerNodes {
		cmd := command.Worker(ctx, controlPlaneNodes)
		fmt.Printf("Installing worker k3s on %s\n", node.Name)
		ssh.Execute(ctx, client, node, cmd)
	}

	gateway.Delete(ctx)
}
