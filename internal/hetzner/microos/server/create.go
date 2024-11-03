// Package server contains the functionality for managing Hetzner cloud microOS servers
package server

import (
	"errors"
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/utils/logger"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

const (
	// ARMInstanceType CAX11 is the server type for ARM architecture
	ARMInstanceType = config.CAX11
	// X86InstanceType CX22 is the server type for x86 architecture
	X86InstanceType = config.CX22
	// LinuxImage is the name of the Ubuntu 24.04 image
	LinuxImage = "ubuntu-24.04"
)

// Create creates the Hetzner cloud microOS server
func Create(
	ctx *cluster.Cluster,
	architecture hcloud.Architecture,
	sshKey *hcloud.SSHKey,
	location config.Location,
) (*hcloud.Server, error) {
	server, err := Get(ctx, architecture)
	if server != nil && err == nil {
		return server, nil
	}
	return create(ctx, architecture, sshKey, location)
}

func create(
	ctx *cluster.Cluster,
	architecture hcloud.Architecture,
	sshKey *hcloud.SSHKey,
	loc config.Location,
) (*hcloud.Server, error) {
	name := getName(ctx, architecture)

	l := logger.New(nil, logger.Server, logger.Create, name)
	defer l.LogEvents()

	instance := ARMInstanceType
	if architecture == hcloud.ArchitectureX86 {
		instance = X86InstanceType
	}

	res, _, err := ctx.CloudClient.Server.Create(ctx.Context, hcloud.ServerCreateOpts{
		Name:       name,
		Image:      &hcloud.Image{Name: LinuxImage},
		ServerType: &hcloud.ServerType{Name: string(instance)},
		Location:   &hcloud.Location{Name: string(loc)},
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
		l.AddEvent(logger.Failure, err)
		return nil, err
	}
	if res.Server == nil {
		err = errors.New("server is nil")
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	if err := ctx.CloudClient.Action.WaitFor(ctx.Context, res.Action); err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}
	if err := ctx.CloudClient.Action.WaitFor(ctx.Context, res.NextActions...); err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	l.AddEvent(logger.Success)
	return res.Server, nil
}
