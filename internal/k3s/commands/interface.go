package commands

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/utils/ssh"
	"strings"
)

func GetNetworkInterfaceName(
	ctx *cluster.Cluster,
	proxy *hcloud.Server,
	remote *hcloud.Server,
) (string, error) {
	cmd := "ip -o link show | awk '$2 != \"lo:\" {print $2}' | sed 's/://g' | head -n 1"
	res, err := ssh.ExecuteViaProxy(ctx, proxy, remote, cmd)
	if err != nil {
		return "", err
	}
	res = strings.TrimSpace(res)
	return res, nil
}
