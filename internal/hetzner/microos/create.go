package microos

import (
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/hetzner/microos/image"
	"h3s/internal/hetzner/microos/server"
	"h3s/internal/hetzner/sshkey"
	"h3s/internal/utils/logger"
	"sync"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func Create(ctx *cluster.Cluster) error {
	sshKey, err := sshkey.Get(ctx)
	if err != nil {
		return err
	}
	architectures := config.GetArchitectures(ctx.Config)

	var wg sync.WaitGroup

	if architectures.ARM {
		wg.Add(1)
		go func() {
			defer wg.Done()
			create(ctx, sshKey, hcloud.ArchitectureARM)
		}()
	}

	if architectures.X86 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			create(ctx, sshKey, hcloud.ArchitectureX86)
		}()
	}

	wg.Wait()
	return nil
}

func create(
	ctx *cluster.Cluster,
	sshKey *hcloud.SSHKey,
	architecture hcloud.Architecture,
) (*hcloud.Image, error) {
	l := logger.New(nil, logger.Image, logger.Create, string(architecture))
	defer l.LogEvents()

	// Check if image already exists - and return if it does
	img, err := image.Get(ctx, architecture)
	if err == nil && img != nil {
		l.AddEvent(logger.Success, "image already exists")
		return img, nil
	}

	// Cleanup on success AND failure
	defer server.Delete(ctx, architecture)

	// Create server that will be used to create the image
	s, err := server.Create(ctx, architecture, sshKey, ctx.Config.ControlPlane.Pool.Location)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	// Start in rescue mode - Required for image creation
	if _, err := server.RescueMode(ctx, sshKey, s); err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	// Setup Image - Download Image & Install Dependencies etc.
	if err := image.Provision(ctx, architecture, s); err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	// Create snapshot/image
	if err := server.Shutdown(ctx, s); err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	// Make the snapshot/image
	return image.Create(ctx, s, architecture)
}
