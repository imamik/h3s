package ssh

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/clustercontext"
	"h3s/internal/resources/gateway"
	"h3s/internal/resources/pool/node"
	"h3s/internal/resources/server"
	"h3s/internal/utils/ssh"
)

func SSH(ctx clustercontext.ClusterContext, cmd string) (string, error) {
	gate, err := gateway.Get(ctx)
	if err != nil {
		return "", err
	}

	nodes := server.GetAll(ctx)
	var controlPlane *hcloud.Server

	for _, n := range nodes {
		if node.IsControlPlaneNode(n) {
			controlPlane = n
			break
		}
	}

	res, err := ssh.ExecuteViaProxy(ctx, gate, controlPlane, cmd)
	if err != nil {
		return "", err
	}

	return res, nil
}
