// Package common contains the common utility functionality
package common

import (
	"h3s/internal/cluster"
	"h3s/internal/hetzner/gateway"
	"h3s/internal/hetzner/server"
	"h3s/internal/utils/ssh"
)

// SSH executes a command on the first control plane node via the gateway
func SSH(ctx *cluster.Cluster, cmd string) (string, error) {
	// Get the gateway
	gate, err := gateway.Get(ctx)
	if err != nil {
		return "", err
	}

	// Get the first control plane node
	nodes, err := server.GetAll(ctx)
	if err != nil {
		return "", err
	}
	firstControlPlane := nodes.ControlPlane[0]

	// Execute the command on the first control plane node
	res, err := ssh.ExecuteViaProxy(ctx.Config.SSHKeyPaths.PrivateKeyPath, gate, firstControlPlane, cmd)
	if err != nil {
		return "", err
	}

	return res, nil
}
