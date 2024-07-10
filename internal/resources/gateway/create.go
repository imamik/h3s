package gateway

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/config"
	"hcloud-k3s-cli/internal/resources/microos/image"
	"hcloud-k3s-cli/internal/resources/network"
	"hcloud-k3s-cli/internal/resources/sshkey"
	"hcloud-k3s-cli/internal/utils/ip"
	"hcloud-k3s-cli/internal/utils/logger"
	"hcloud-k3s-cli/internal/utils/ssh"
)

func Create(ctx clustercontext.ClusterContext) *hcloud.Server {
	net := network.Get(ctx)
	sshKey := sshkey.Get(ctx)

	img, _ := image.Get(ctx, hcloud.ArchitectureARM)
	proxy := createServer(ctx, sshKey, net, img)
	configureGateway(ctx, proxy)
	setupGatewayRoute(ctx, net, proxy)

	return proxy
}

func configureGateway(ctx clustercontext.ClusterContext, proxy *hcloud.Server) {
	ssh.ExecuteWithSsh(ctx, proxy, `
echo 1 > /proc/sys/net/ipv4/ip_forward
iptables -t nat -A POSTROUTING -s '10.0.0.0/16' -o eth0 -j MASQUERADE
	`)
}

func setupGatewayRoute(ctx clustercontext.ClusterContext, net *hcloud.Network, proxy *hcloud.Server) {
	ctx.Client.Network.AddRoute(ctx.Context, net, hcloud.NetworkAddRouteOpts{
		Route: hcloud.NetworkRoute{
			Destination: ip.GetIpRange("0.0.0.0/0"),
			Gateway:     proxy.PrivateNet[0].IP,
		},
	})
}

func createServer(
	ctx clustercontext.ClusterContext,
	sshKey *hcloud.SSHKey,
	network *hcloud.Network,
	image *hcloud.Image,
) *hcloud.Server {
	server, err := getServer(ctx)
	if err == nil && server != nil {
		return server
	}

	name := getName(ctx)
	logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Initialized)

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
			"is_gateway": "true",
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

	server, err = getServer(ctx)
	if err != nil {
		logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Success)

	return server
}