// Package k3s contains the functionality for installing k3s
package k3s

import (
	"h3s/internal/cluster"
	"h3s/internal/hetzner/gateway"
	"h3s/internal/hetzner/loadbalancers"
	"h3s/internal/hetzner/server"
	"h3s/internal/k3s/config"
	"h3s/internal/utils/ip"
	"h3s/internal/utils/ssh"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// getNetworkInterface calls the ip command on the remote server to get the network interface name
func getNetworkInterface(ctx *cluster.Cluster, proxy, remote *hcloud.Server) (string, error) {
	cmd := "ip -o link show | awk '$2 != \"lo:\" {print $2}' | sed 's/://g' | head -n 1"
	return ssh.ExecuteViaProxy(ctx.Config.SSHKeyPaths.PrivateKeyPath, proxy, remote, cmd)
}

// getServer returns the server address
func getServer(firstControlPlane *hcloud.Server) string {
	return "https://" + ip.Private(firstControlPlane).String() + ":6443"
}

// getTLSSan returns the SANs (Subject Alternative Names) for the TLS certificate of the k3s server
func getTLSSan(domain string, lb *hcloud.LoadBalancer, controlPlaneNodes []*hcloud.Server) []string {
	tlsSan := []string{
		"127.0.0.1",
		"localhost",
		"kubernetes",
		"kubernetes.default",
		"kubernetes.default.svc",
		domain,
		domain + ".local",
		"k3s." + domain,
		lb.PublicNet.IPv4.IP.String(),
		lb.PublicNet.IPv6.IP.String(),
	}

	for _, privateNet := range lb.PrivateNet {
		tlsSan = append(tlsSan, privateNet.IP.String())
	}

	for _, node := range controlPlaneNodes {
		tlsSan = append(tlsSan, ip.Private(node).String())
	}

	return tlsSan
}

// Install installs k3s to the cluster
func Install(ctx *cluster.Cluster) error {
	// Get load balancer
	lb, err := loadbalancers.Get(ctx)
	if err != nil {
		return err
	}

	// Get gateway
	gate, err := gateway.Get(ctx)
	if err != nil {
		return err
	}

	// Get all nodes
	nodes, err := server.GetAll(ctx)
	if err != nil {
		return err
	}
	main := nodes.ControlPlane[0]
	serverAddress := getServer(main)
	tlsSAN := getTLSSan(ctx.Config.Domain, lb, nodes.ControlPlane)

	install := func(n *hcloud.Server, isControlPlane, isMain bool) error {
		networkInterface, err := getNetworkInterface(ctx, gate, n)
		if err != nil {
			return err
		}

		// Get the private IP address of the node
		nodeIP := ip.Private(n)
		cmd, err := config.Command(config.CommandConfig{
			IsMain:                isMain,
			IsControlPlane:        isControlPlane,
			K3sToken:              ctx.Credentials.K3sToken,
			Server:                serverAddress,
			Domain:                ctx.Config.Domain,
			TLSSAN:                tlsSAN,
			ControlPlanesAsWorker: ctx.Config.ControlPlane.AsWorkerPool,
			NodeName:              n.Name,
			NodeIP:                nodeIP.String(),
			NetworkInterface:      networkInterface,
			K3sVersion:            ctx.Config.K3sVersion,
		})
		if err != nil {
			return err
		}

		_, err = ssh.ExecuteViaProxy(ctx.Config.SSHKeyPaths.PrivateKeyPath, gate, n, cmd)
		return err
	}

	for _, n := range nodes.ControlPlane {
		if err := install(n, true, n == main); err != nil {
			return err
		}
	}

	for _, n := range nodes.Worker {
		if err := install(n, false, false); err != nil {
			return err
		}
	}

	return nil
}
