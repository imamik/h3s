package ssh

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/melbahja/goph"
	"hcloud-k3s-cli/pkg/clustercontext"
)

func Execute(
	ctx clustercontext.ClusterContext,
	client *goph.Client,
	server *hcloud.Server,
	command string,
) {
	privateKeyPath := ctx.Config.SSHKeyPaths.PrivateKeyPath
	serverIp := server.PrivateNet[0].IP.String()

	command = "ssh -o StrictHostKeyChecking=no -i " + privateKeyPath + " root@" + serverIp + " " + command
	run(client, command)
}
