package install

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/melbahja/goph"
	"hcloud-k3s-cli/pkg/clustercontext"
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

	gate := gateway.Create(ctx)
	client := ssh.Client(ctx, gate)

	defer func(client *goph.Client) {
		err := client.Close()
		if err != nil {
			return
		}
	}(client)

	for _, node := range nodes {
		installOnNode(ctx, client, node)
	}

	gateway.Delete(ctx)
}

func installOnNode(
	ctx clustercontext.ClusterContext,
	client *goph.Client,
	node *hcloud.Server,
) {
	version := ctx.Config.K3sVersion
	command := fmt.Sprintf("curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=%s sh -", version)

	//isControlPlane := node.Labels["is_control_plane"] == "true"
	//isWorker := node.Labels["is_worker"] == "true"

	fmt.Printf("Installing k3s on %s\n", node.Name)
	ssh.Execute(ctx, client, node, command)
}
