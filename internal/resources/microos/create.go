package microos

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/microos/image"
	"hcloud-k3s-cli/internal/resources/microos/server"
	"hcloud-k3s-cli/internal/resources/network"
	"hcloud-k3s-cli/internal/resources/sshkey"
)

func Create(
	ctx clustercontext.ClusterContext,
	architecture hcloud.Architecture,
) *hcloud.Image {
	img, err := image.Get(ctx, architecture)
	if err == nil && img != nil {
		return img
	}

	// Prepare server to create image from
	sshKey := sshkey.Create(ctx)
	net := network.Create(ctx)
	s := server.Create(ctx, architecture, sshKey, net, ctx.Config.ControlPlane.Pool.Location)
	rootPassword := server.RescueMode(ctx, sshKey, s)

	// Setup Image - Download Image & Install Dependencies etc.
	image.Provision(architecture, s, rootPassword)

	// Create snapshot/image
	server.Shutdown(ctx, s)
	img = image.Create(ctx, s, architecture)

	// Clean up
	server.Delete(ctx, architecture)

	return img
}
