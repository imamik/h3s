package ip

import (
	"hcloud-k3s-cli/internal/utils/logger"
	"net"
)

func GetIpRange(s string) *net.IPNet {
	_, ipRange, err := net.ParseCIDR(s)
	if err != nil {
		logger.LogError("Invalid IP Range", err)
	}
	return ipRange
}
