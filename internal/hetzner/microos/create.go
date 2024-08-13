package microos

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/hetzner/microos/image"
	"h3s/internal/hetzner/microos/server"
	"h3s/internal/hetzner/sshkey"
	"sync"
)

func Create(ctx *cluster.Cluster) {
	sshKey := sshkey.Get(ctx)
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
}

func create(
	ctx *cluster.Cluster,
	sshKey *hcloud.SSHKey,
	architecture hcloud.Architecture,
) *hcloud.Image {
	img, err := image.Get(ctx, architecture)
	if err == nil && img != nil {
		return img
	}

	s := server.Create(ctx, architecture, sshKey, ctx.Config.ControlPlane.Pool.Location)
	defer server.Delete(ctx, architecture) // Cleanup on success AND failure
	server.RescueMode(ctx, sshKey, s)

	// Setup Image - Download Image & Install Dependencies etc.
	image.Provision(ctx, architecture, s)

	// Create snapshot/image
	server.Shutdown(ctx, s)
	img = image.Create(ctx, s, architecture)

	return img
}
