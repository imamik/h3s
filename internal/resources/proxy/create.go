package proxy

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/config"
	"hcloud-k3s-cli/internal/resources/network"
	"hcloud-k3s-cli/internal/resources/sshkey"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Create(ctx clustercontext.ClusterContext) *hcloud.Server {
	net := network.Get(ctx)
	sshKey := sshkey.Get(ctx)
	proxy := createServer(ctx, sshKey, net)
	return proxy
}

func createServer(
	ctx clustercontext.ClusterContext,
	sshKey *hcloud.SSHKey,
	network *hcloud.Network,
) *hcloud.Server {
	server, err := getServer(ctx)
	if err == nil && server != nil {
		return server
	}

	name := getName(ctx)
	logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Initialized)

	image := &hcloud.Image{Name: "ubuntu-24.04"}
	serverType := &hcloud.ServerType{Name: string(config.CAX11)}
	location := &hcloud.Location{Name: string(ctx.Config.ControlPlane.Pool.Location)}
	publicNet := &hcloud.ServerCreatePublicNet{
		EnableIPv4: true,
		EnableIPv6: true,
	}
	networks := []*hcloud.Network{network}
	sshKeys := []*hcloud.SSHKey{sshKey}

	res, _, err := ctx.Client.Server.Create(ctx.Context, hcloud.ServerCreateOpts{
		Name:       name,
		ServerType: serverType,
		Image:      image,
		Location:   location,
		Networks:   networks,
		SSHKeys:    sshKeys,
		PublicNet:  publicNet,
		Labels: ctx.GetLabels(map[string]string{
			"is_proxy": "true",
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
