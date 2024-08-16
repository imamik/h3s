package ip

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/utils/logger"
	"net"
	"net/netip"
)

// FirstAvailable returns the first available IP address of a server
// It will return the public IPv4 address if available,
// otherwise the public IPv6 address if available
// and fallback to the first private IP address
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

// Private returns the first private IP address of a server
func Private(server *hcloud.Server) string {
	if len(server.PrivateNet) > 0 {
		return server.PrivateNet[0].IP.String()
	}
	return ""
}

// GetIpRange returns the IP range of a CIDR string
func GetIpRange(s string) *net.IPNet {
	_, ipRange, err := net.ParseCIDR(s)
	if err != nil {
		logger.LogError("Invalid IP Range", err)
	}
	return ipRange
}
