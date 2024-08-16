package k8s

import (
	"fmt"
	"h3s/internal/cluster"
	"h3s/internal/hetzner/gateway"
	"h3s/internal/hetzner/server"
	"h3s/internal/utils/kubectl"
	"h3s/internal/utils/ssh"
)

func SetServer(clr *cluster.Cluster) error {
	gatewayNode, _ := gateway.GetIfNeeded(clr)
	nodes, err := server.GetAll(clr)
	if err != nil {
		return err
	}
	firstControlPlane := nodes.ControlPlane[0]
	serverArg := fmt.Sprintf("--server=https://k3s.%s", clr.Config.Domain)
	cmd, err := kubectl.New("config", "set-cluster", "default", serverArg).String()
	if err != nil {
		return err
	}
	_, err = ssh.ExecuteViaProxy(clr.Config.SSHKeyPaths.PrivateKeyPath, gatewayNode, firstControlPlane, cmd)
	return err
}
