package common

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/hetzner/gateway"
	"h3s/internal/hetzner/pool/node"
	"h3s/internal/hetzner/server"
	"h3s/internal/utils/ssh"
)

func SSH(ctx *cluster.Cluster, cmd string) (string, error) {
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
