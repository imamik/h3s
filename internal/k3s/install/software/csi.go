package software

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
)

func InstallHetznerCSI(
	ctx clustercontext.ClusterContext,
	network *hcloud.Network,
	proxy *hcloud.Server,
	remote *hcloud.Server,
) {
	ApplyYaml(ctx, proxy, remote, "https://raw.githubusercontent.com/hetznercloud/csi-driver/v2.6.0/deploy/kubernetes/hcloud-csi.yml")
}
