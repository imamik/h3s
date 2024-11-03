// Package ping contains the functionality for pinging servers
package ping

import (
	"fmt"
	"h3s/internal/utils/execute"
	"h3s/internal/utils/ip"
	"h3s/internal/utils/logger"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Ping pings the server every 5 seconds until it is available.
func Ping(server *hcloud.Server, timeout time.Duration) {
	l := logger.New(nil, logger.Server, "Ping", server.Name)
	defer l.LogEvents()
	ipAddress := ip.FirstAvailable(server)

	for {
		if isServerAvailable(ipAddress) {
			l.AddEvent(logger.Success)
			break
		}
		l.AddEvent(logger.Failure, fmt.Sprintf("Server %s is not available", server.Name), fmt.Sprintf("Retry in %s", timeout))
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
