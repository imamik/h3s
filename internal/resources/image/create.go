package image

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/config"
	"hcloud-k3s-cli/internal/resources/image/server"
	"hcloud-k3s-cli/internal/resources/sshkey"
)

func Create(
	ctx clustercontext.ClusterContext,
	architecture hcloud.Architecture,
	l config.Location,
) *hcloud.Image {
	img, err := Get(ctx, architecture)
	if err == nil && img != nil {
		return img
	}
	sshKey := sshkey.Create(ctx)
	s := server.Create(ctx, sshKey, architecture, l)
	img = createImage(ctx, s)
	server.Delete(ctx, architecture)
	return img
}

func createImage(
	ctx clustercontext.ClusterContext,
	server *hcloud.Server,
) *hcloud.Image {
	// log.Fatalf("createImage not implemented")
	return nil
}
