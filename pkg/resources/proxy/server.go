package proxy

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/clustercontext"
	"hcloud-k3s-cli/pkg/config"
	"hcloud-k3s-cli/pkg/utils/logger"
)

func getName(ctx clustercontext.ClusterContext) string {
	return ctx.GetName("proxy")
}

func getServer(ctx clustercontext.ClusterContext) (*hcloud.Server, error) {
	name := getName(ctx)
	logger.LogResourceEvent(logger.Server, logger.Get, name, logger.Initialized)

	server, _, err := ctx.Client.Server.GetByName(ctx.Context, name)
	if err != nil {
		logger.LogResourceEvent(logger.Server, logger.Get, name, logger.Failure, err)
		return nil, err
	}

	logger.LogResourceEvent(logger.Server, logger.Get, name, logger.Success)
	return server, nil
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

	logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Success)
	return res.Server
}

func deleteServer(ctx clustercontext.ClusterContext) {
	server, err := getServer(ctx)
	if err != nil || server == nil {
		return
	}

	logger.LogResourceEvent(logger.Server, logger.Delete, server.Name, logger.Initialized)

	_, _, err = ctx.Client.Server.DeleteWithResult(ctx.Context, server)
	if err != nil {
		logger.LogResourceEvent(logger.Server, logger.Delete, server.Name, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.Server, logger.Delete, server.Name, logger.Success)
}
