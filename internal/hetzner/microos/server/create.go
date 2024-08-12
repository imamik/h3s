package server

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/utils/logger"
)

const (
	// ARMInstanceType CAX11 is the server type for ARM architecture
	ARMInstanceType = config.CAX11
	// X86InstanceType CX22 is the server type for x86 architecture
	X86InstanceType = config.CX22
	LinuxImage      = "ubuntu-24.04"
)

func Create(
	ctx *cluster.Cluster,
	architecture hcloud.Architecture,
	sshKey *hcloud.SSHKey,
	location config.Location,
) *hcloud.Server {
	server := Get(ctx, architecture)
	if server != nil {
		return server
	}
	return create(ctx, architecture, sshKey, location)
}

func create(
	ctx *cluster.Cluster,
	architecture hcloud.Architecture,
	sshKey *hcloud.SSHKey,
	l config.Location,
) *hcloud.Server {
	name := getName(ctx, architecture)
	logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Initialized)

	instance := ARMInstanceType
	if architecture == hcloud.ArchitectureX86 {
		instance = X86InstanceType
	}

	res, _, err := ctx.CloudClient.Server.Create(ctx.Context, hcloud.ServerCreateOpts{
		Name:       name,
		Image:      &hcloud.Image{Name: LinuxImage},
		ServerType: &hcloud.ServerType{Name: string(instance)},
		Location:   &hcloud.Location{Name: string(l)},
		SSHKeys:    []*hcloud.SSHKey{sshKey},
		PublicNet: &hcloud.ServerCreatePublicNet{
			EnableIPv4: true,
			EnableIPv6: true,
		},
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

	if err := ctx.CloudClient.Action.WaitFor(ctx.Context, res.Action); err != nil {
		logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Failure, err)
	}
	if err := ctx.CloudClient.Action.WaitFor(ctx.Context, res.NextActions...); err != nil {
		logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Success)
	return res.Server
}
