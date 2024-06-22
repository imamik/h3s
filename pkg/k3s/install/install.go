package install

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/clustercontext"
	"hcloud-k3s-cli/pkg/utils/ssh"
)

func Install(
	ctx clustercontext.ClusterContext,
	server *hcloud.Server,
) {
	version := ctx.Config.K3sVersion
	command := fmt.Sprintf("curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=%s sh -", version)

	isControlPlane := server.Labels["is_control_plane"] == "true"
	isWorker := server.Labels["is_worker"] == "true"

	fmt.Printf("Installing k3s on %s\n", server.Name)
	fmt.Printf("Control Plane: %t\n", isControlPlane)
	fmt.Printf("Worker: %t\n", isWorker)

	ssh.Execute(ctx, server, command)
}
