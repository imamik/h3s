package server

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/config"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Create(
	ctx clustercontext.ClusterContext,
	sshKey *hcloud.SSHKey,
	architecture hcloud.Architecture,
	l config.Location,
) *hcloud.Server {
	name := getName(ctx, architecture)
	logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Initialized)

	image := &hcloud.Image{Name: "ubuntu-24.04"}

	instance := config.CAX11
	if architecture == hcloud.ArchitectureX86 {
		instance = config.CX22
	}
	serverType := &hcloud.ServerType{Name: string(instance)}
	location := &hcloud.Location{Name: string(l)}
	publicNet := &hcloud.ServerCreatePublicNet{
		EnableIPv4: true,
		EnableIPv6: true,
	}
	sshKeys := []*hcloud.SSHKey{sshKey}

	res, _, err := ctx.Client.Server.Create(ctx.Context, hcloud.ServerCreateOpts{
		Name:       name,
		ServerType: serverType,
		Image:      image,
		Location:   location,
		SSHKeys:    sshKeys,
		PublicNet:  publicNet,
		Labels: ctx.GetLabels(map[string]string{
			"is_image_creator": "true",
			"architecture":     string(architecture),
		}),
	})
	if err != nil {
		logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Failure, err)
	}
	if res.Server == nil {
		logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Failure, "Empty Response")
	}

	logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Success)
	return res.Server
}
