package microos

import (
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/hetzner/microos/image"
	"h3s/internal/hetzner/microos/server"
	"h3s/internal/utils/logger"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Delete deletes the MicroOS image and server for the given cluster.
func Delete(ctx *cluster.Cluster) {
	architectures := config.GetArchitectures(ctx.Config)

	if architectures.ARM {
		if err := del(ctx, hcloud.ArchitectureARM); err != nil {
			logger.New(nil, logger.Image, logger.Delete, "ARM").AddEvent(logger.Failure, err)
		}
	}

	if architectures.X86 {
		if err := del(ctx, hcloud.ArchitectureX86); err != nil {
			logger.New(nil, logger.Image, logger.Delete, "X86").AddEvent(logger.Failure, err)
		}
	}
}

func del(ctx *cluster.Cluster, architecture hcloud.Architecture) error {
	err := image.Delete(ctx, architecture)
	if err != nil {
		return err
	}

	err = server.Delete(ctx, architecture)
	if err != nil {
		return err
	}

	return nil
}
