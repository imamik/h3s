package microos

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/clustercontext"
	"h3s/internal/config"
	"h3s/internal/resources/microos/image"
	"h3s/internal/resources/microos/server"
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
