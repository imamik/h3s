package ssh

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/clustercontext"
)

func Execute(
	ctx clustercontext.ClusterContext,
	server *hcloud.Server,
	command string,
) {
	ip := server.PublicNet.IPv4.IP.String()
	privateKeyPath := ctx.Config.SSHKeyPaths.PrivateKeyPath
	client := getClient(ip, privateKeyPath, "")
	run(client, command)
}
