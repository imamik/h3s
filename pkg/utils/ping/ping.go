package ping

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"os/exec"
	"time"
)

// Ping pings the server every 5 seconds until it is available.
func Ping(server *hcloud.Server) {
	ip := server.PublicNet.IPv4.IP.String()

	for {
		if isServerAvailable(ip) {
			fmt.Printf("Server %s is available.\n", ip)
			break
		}
		fmt.Printf("Server %s (%s) is not available. Retrying in 10 seconds...\n", server.Name, ip)
		time.Sleep(10 * time.Second)
	}

	time.Sleep(10 * time.Second)
}

// isServerAvailable uses the system ping command to check if the server is available.
func isServerAvailable(ip string) bool {
	cmd := exec.Command("ping", "-c", "1", ip)
	err := cmd.Run()
	return err == nil
}