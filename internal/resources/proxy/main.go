package proxy

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/network"
	"hcloud-k3s-cli/internal/resources/sshkey"
	"hcloud-k3s-cli/internal/utils/ping"
)

func Create(ctx clustercontext.ClusterContext) *hcloud.Server {
	net := network.Get(ctx)
	sshKey := sshkey.Get(ctx)
	gateway := createServer(ctx, sshKey, net)
	ping.Ping(gateway)
	return gateway
}

func Delete(ctx clustercontext.ClusterContext) {
	deleteServer(ctx)
}
