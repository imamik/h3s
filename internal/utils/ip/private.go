package ip

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func Private(server *hcloud.Server) string {
	if len(server.PrivateNet) > 0 {
		return server.PrivateNet[0].IP.String()
	}
	return ""
}
