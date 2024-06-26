package microos

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/microos/image"
	"hcloud-k3s-cli/internal/resources/microos/server"
)

func Delete(ctx clustercontext.ClusterContext, architecture hcloud.Architecture) {
	image.Delete(ctx, architecture)
	server.Delete(ctx, architecture)
}
