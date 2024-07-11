package ip

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"net/netip"
)

func FirstAvailable(server *hcloud.Server) string {
	switch {
	case !server.PublicNet.IPv4.IsUnspecified():
		return server.PublicNet.IPv4.IP.String()
	case !server.PublicNet.IPv6.IsUnspecified():
		network, ok := netip.AddrFromSlice(server.PublicNet.IPv6.IP)
		if ok {
			return fmt.Sprintf("[%s]", network.Next().String())
		}
	default:
		return Private(server)
	}
	return ""
}
