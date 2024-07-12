package commands

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/utils/ip"
)

func getServer(
	lb *hcloud.LoadBalancer,
	node *hcloud.Server,
) string {
	address := ""
	if lb == nil {
		address = ip.Private(node)
	} else if len(lb.PrivateNet) > 0 {
		address = lb.PrivateNet[0].IP.String()
	} else {
		address = lb.PublicNet.IPv4.IP.String()
	}
	address = ip.Private(node)
	return "https://" + address + ":6443"
}

func getTlsSan(
	lb *hcloud.LoadBalancer,
	controlPlaneNodes []*hcloud.Server,
) []string {
	var tlsSan []string

	if lb == nil {
		for _, node := range controlPlaneNodes {
			tlsSan = append(tlsSan, ip.FirstAvailable(node))
		}
	} else {
		tlsSan = append(tlsSan, lb.PublicNet.IPv4.IP.String())
		tlsSan = append(tlsSan, lb.PublicNet.IPv6.IP.String())
		for _, privateNet := range lb.PrivateNet {
			tlsSan = append(tlsSan, privateNet.IP.String())
		}
	}

	return tlsSan
}
