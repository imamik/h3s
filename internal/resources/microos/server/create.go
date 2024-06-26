package server

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/config"
	"hcloud-k3s-cli/internal/utils/logger"
)

const (
	// ARMInstanceType CAX11 is the server type for ARM architecture
	ARMInstanceType = config.CAX11
	// X86InstanceType CX22 is the server type for x86 architecture
	X86InstanceType = config.CX22
	LinuxImage      = "ubuntu-24.04"
)

func Create(
	ctx clustercontext.ClusterContext,
	sshKey *hcloud.SSHKey,
	architecture hcloud.Architecture,
	l config.Location,
) *hcloud.Server {
	name := getName(ctx, architecture)
	logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Initialized)

	image := &hcloud.Image{Name: LinuxImage}

	instance := ARMInstanceType
	if architecture == hcloud.ArchitectureX86 {
		instance = X86InstanceType
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

	if err := ctx.Client.Action.WaitFor(ctx.Context, res.Action); err != nil {
		logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Failure, err)
	}
	if err := ctx.Client.Action.WaitFor(ctx.Context, res.NextActions...); err != nil {
		logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Success)
	return res.Server
}
