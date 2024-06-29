package microos

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/config"
	"hcloud-k3s-cli/internal/resources/microos/image"
	"hcloud-k3s-cli/internal/resources/microos/server"
)

func Delete(ctx clustercontext.ClusterContext) {
	architectures := config.GetArchitectures(ctx.Config)

	if architectures.ARM {
		del(ctx, hcloud.ArchitectureARM)
	}

	if architectures.X86 {
		del(ctx, hcloud.ArchitectureX86)
	}
}

func del(ctx clustercontext.ClusterContext, architecture hcloud.Architecture) {
	image.Delete(ctx, architecture)
	server.Delete(ctx, architecture)
}
