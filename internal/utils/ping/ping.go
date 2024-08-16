package ping

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/utils/execute"
	"h3s/internal/utils/ip"
	"h3s/internal/utils/logger"
	"time"
)

// Ping pings the server every 5 seconds until it is available.
func Ping(server *hcloud.Server, timeout time.Duration) {
	ipAddress := ip.FirstAvailable(server)

	for {
		if isServerAvailable(ipAddress) {
			logger.LogResourceEvent(logger.Server, "Available", server.Name, logger.Success)
			break
		}
		logger.LogResourceEvent(logger.Server, "Not Available", server.Name, logger.Failure, fmt.Sprintf("Retry in %s", timeout))
		time.Sleep(timeout)
	}

	time.Sleep(timeout)
}

// isServerAvailable uses the system ping command to check if the server is available.
func isServerAvailable(ip string) bool {
	cmd := fmt.Sprintf("ping -c 1 %s", ip)
	_, err := execute.Local(cmd)
	return err == nil
}
