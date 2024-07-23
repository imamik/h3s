package commands

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/ip"
)

func getServer(firstControlPlane *hcloud.Server) string {
	return "https://" + ip.Private(firstControlPlane) + ":6443"
}

func getTlsSan(
	ctx clustercontext.ClusterContext,
	lb *hcloud.LoadBalancer,
	controlPlaneNodes []*hcloud.Server,
) []string {
	tlsSan := []string{
		"127.0.0.1",
		"localhost",
		"kubernetes",
		"kubernetes.default",
		"kubernetes.default.svc",
		ctx.Config.Domain,
		ctx.Config.Domain + ".local",
		"k3s." + ctx.Config.Domain,
		lb.PublicNet.IPv4.IP.String(),
		lb.PublicNet.IPv6.IP.String(),
	}

	for _, privateNet := range lb.PrivateNet {
		tlsSan = append(tlsSan, privateNet.IP.String())
	}

	for _, node := range controlPlaneNodes {
		tlsSan = append(tlsSan, node.PrivateNet[0].IP.String())
	}

	return tlsSan
}
