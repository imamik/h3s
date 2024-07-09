package software

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
)

func Install(
	ctx clustercontext.ClusterContext,
	net *hcloud.Network,
	lb *hcloud.LoadBalancer,
	proxyServer *hcloud.Server,
	remote *hcloud.Server,
) {
	InstallHetznerCCM(ctx, net, proxyServer, remote)
	InstallHetznerCSI(ctx, net, proxyServer, remote)
	InstallTraefik(ctx, lb, proxyServer, remote)
}
