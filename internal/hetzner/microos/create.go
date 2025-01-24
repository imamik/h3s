// Package microos contains the functionality for managing Hetzner cloud microOS images
package microos

import (
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/hetzner/microos/image"
	"h3s/internal/hetzner/microos/server"
	"h3s/internal/utils/logger"
	"sync"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Create creates the Hetzner cloud microOS image
func Create(ctx *cluster.Cluster, sshKey *hcloud.SSHKey) (*ImageInArchitecture, error) {
	architectures := config.GetArchitectures(ctx.Config)
	images := &ImageInArchitecture{}

	var wg sync.WaitGroup

	if architectures.ARM {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if img, err := create(ctx, sshKey, hcloud.ArchitectureARM); err != nil {
				logger.New(nil, logger.Image, logger.Create, "ARM").AddEvent(logger.Failure, err)
			} else {
				images.ARM = img
			}
		}()
	}

	if architectures.X86 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if img, err := create(ctx, sshKey, hcloud.ArchitectureX86); err != nil {
				logger.New(nil, logger.Image, logger.Create, "X86").AddEvent(logger.Failure, err)
			} else {
				images.X86 = img
			}
		}()
	}

	wg.Wait()
	return images, nil
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
	defer func() {
		if deleteErr := server.Delete(ctx, architecture); deleteErr != nil {
			l.AddEvent(logger.Failure, deleteErr) // Log the error if deletion fails
		}
	}()

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
