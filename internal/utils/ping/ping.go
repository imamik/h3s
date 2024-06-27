package ping

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"os/exec"
	"time"
)

// Ping pings the server every 5 seconds until it is available.
func Ping(server *hcloud.Server, timeout time.Duration) {
	ip := server.PublicNet.IPv4.IP.String()

	for {
		if isServerAvailable(ip) {
			fmt.Printf("Server %s is available.\n", ip)
			break
		}
		fmt.Printf("Server %s (%s) is not available. Retrying in %s seconds...\n", server.Name, ip, timeout)
		time.Sleep(timeout * time.Second)
	}

	time.Sleep(timeout * time.Second)
}

// isServerAvailable uses the system ping command to check if the server is available.
func isServerAvailable(ip string) bool {
	cmd := exec.Command("ping", "-c", "1", ip)
	err := cmd.Run()
	return err == nil
}
