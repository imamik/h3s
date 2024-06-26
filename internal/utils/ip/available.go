package ip

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"net/netip"
)

func FirstAvailableIP(server *hcloud.Server) string {
	switch {
	case !server.PublicNet.IPv4.IsUnspecified():
		return server.PublicNet.IPv4.IP.String()
	case !server.PublicNet.IPv6.IsUnspecified():
		network, ok := netip.AddrFromSlice(server.PublicNet.IPv6.IP)
		if ok {
			return network.Next().String()
		}
	case len(server.PrivateNet) > 0:
		return server.PrivateNet[0].IP.String()
	}
	return ""
}
