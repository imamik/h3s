package kubectl

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/gateway"
	"hcloud-k3s-cli/internal/resources/pool/node"
	"hcloud-k3s-cli/internal/resources/server"
	"hcloud-k3s-cli/internal/utils/ssh"
	"strings"
)

func SSH(ctx clustercontext.ClusterContext, args []string) error {
	gate, err := gateway.Get(ctx)
	if err != nil {
		return err
	}

	nodes := server.GetAll(ctx)
	var controlPlane *hcloud.Server

	for _, n := range nodes {
		if node.IsControlPlaneNode(n) {
			controlPlane = n
			break
		}
	}

	args = append([]string{"kubectl"}, args...)
	cmd := strings.Join(args, " ")

	_, err = ssh.ExecuteViaProxy(ctx, gate, controlPlane, cmd)
	if err != nil {
		return err
	}

	return nil
}
