package k3s

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/hetzner/gateway"
	"h3s/internal/hetzner/loadbalancers"
	"h3s/internal/hetzner/server"
	"h3s/internal/k3s/config"
	"h3s/internal/utils/ip"
	"h3s/internal/utils/ssh"
)

func getNetworkInterface(ctx *cluster.Cluster, proxy *hcloud.Server, remote *hcloud.Server) (string, error) {
	cmd := "ip -o link show | awk '$2 != \"lo:\" {print $2}' | sed 's/://g' | head -n 1"
	return ssh.ExecuteViaProxy(ctx.Config.SSHKeyPaths.PrivateKeyPath, proxy, remote, cmd)
}

func getServer(firstControlPlane *hcloud.Server) string {
	return "https://" + ip.Private(firstControlPlane) + ":6443"
}

func getTlsSan(domain string, lb *hcloud.LoadBalancer, controlPlaneNodes []*hcloud.Server) []string {
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
		tlsSan = append(tlsSan, ip.Private(node))
	}

	return tlsSan
}

func Install(ctx *cluster.Cluster) error {
	lb := loadbalancers.Get(ctx)
	gate, _ := gateway.GetIfNeeded(ctx)

	nodes, err := server.GetAll(ctx)
	if err != nil {
		return err
	}
	main := nodes.ControlPlane[0]
	serverAddress := getServer(main)
	tlsSAN := getTlsSan(ctx.Config.Domain, lb, nodes.ControlPlane)

	install := func(n *hcloud.Server, isControlPlane bool, isMain bool) error {
		networkInterface, err := getNetworkInterface(ctx, gate, n)
		if err != nil {
			return err
		}

		nodeIp := ip.Private(n)
		cmd, err := config.Command(config.CommandConfig{
			IsMain:                isMain,
			IsControlPlane:        isControlPlane,
			K3sToken:              ctx.Credentials.K3sToken,
			Server:                serverAddress,
			Domain:                ctx.Config.Domain,
			TlsSAN:                tlsSAN,
			ControlPlanesAsWorker: ctx.Config.ControlPlane.AsWorkerPool,
			NodeName:              n.Name,
			NodeIp:                nodeIp,
			NetworkInterface:      networkInterface,
			PublicIps:             ctx.Config.PublicIps,
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
