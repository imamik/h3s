package install

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/clustercontext"
	"hcloud-k3s-cli/pkg/k3s/install/command"
	"hcloud-k3s-cli/pkg/resources/proxy"
	"hcloud-k3s-cli/pkg/resources/server"
	"hcloud-k3s-cli/pkg/utils/ssh"
	"log"
	"time"
)

func Install(
	ctx clustercontext.ClusterContext,
) {
	var nodes []*hcloud.Server
	for {
		nodes = server.GetAllByProject(ctx)
		if len(nodes) == 0 {
			log.Fatal("No servers found")
		}

		allNodesHavePrivateIP := true
		for _, node := range nodes {
			if len(node.PrivateNet) < 1 {
				allNodesHavePrivateIP = false
				break
			}
		}

		if allNodesHavePrivateIP {
			break
		}

		fmt.Println("Not all nodes have an assigned private IP. Retrying in 10 seconds...")
		time.Sleep(10 * time.Second)
	}

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
	}

	for _, node := range workerNodes {
		cmd := command.Worker(ctx, controlPlaneNodes)
		fmt.Printf("Installing worker k3s on %s\n", node.Name)
		ssh.Execute(ctx, p, node, cmd)
	}

	proxy.Delete(ctx)
}
