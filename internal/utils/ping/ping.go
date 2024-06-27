package ping

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/utils/logger"
	"os/exec"
	"time"
)

// Ping pings the server every 5 seconds until it is available.
func Ping(server *hcloud.Server, timeout time.Duration) {
	ip := server.PublicNet.IPv4.IP.String()

	for {
		if isServerAvailable(ip) {
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
	cmd := exec.Command("ping", "-c", "1", ip)
	err := cmd.Run()
	return err == nil
}
