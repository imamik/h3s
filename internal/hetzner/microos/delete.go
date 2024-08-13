package microos

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/hetzner/microos/image"
	"h3s/internal/hetzner/microos/server"
)

func Delete(ctx *cluster.Cluster) {
	architectures := config.GetArchitectures(ctx.Config)

	if architectures.ARM {
		del(ctx, hcloud.ArchitectureARM)
	}

	if architectures.X86 {
		del(ctx, hcloud.ArchitectureX86)
	}
}

func del(ctx *cluster.Cluster, architecture hcloud.Architecture) {
	image.Delete(ctx, architecture)
	server.Delete(ctx, architecture)
}
