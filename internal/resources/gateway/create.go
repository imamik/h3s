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
	gateway := createServer(ctx, sshKey, net, img)
	configureGateway(ctx, gateway)
	setupGatewayRoute(ctx, net, gateway)

	return gateway
}

func configureGateway(ctx clustercontext.ClusterContext, gateway *hcloud.Server) {
	ssh.ExecuteWithSsh(ctx, gateway, `
echo 1 > /proc/sys/net/ipv4/ip_forward
iptables -t nat -A POSTROUTING -s '10.0.0.0/16' -o eth0 -j MASQUERADE
	`)
}

func setupGatewayRoute(ctx clustercontext.ClusterContext, net *hcloud.Network, gateway *hcloud.Server) {
	ctx.Client.Network.AddRoute(ctx.Context, net, hcloud.NetworkAddRouteOpts{
		Route: hcloud.NetworkRoute{
			Destination: ip.GetIpRange("0.0.0.0/0"),
			Gateway:     gateway.PrivateNet[0].IP,
		},
	})
}

func createServer(
	ctx clustercontext.ClusterContext,
	sshKey *hcloud.SSHKey,
	network *hcloud.Network,
	image *hcloud.Image,
) *hcloud.Server {
	server, err := Get(ctx)
	if err == nil && server != nil {
		return server
	}

	name := getName(ctx)

	addEvent, logEvents := logger.NewEventLogger(logger.Server, logger.Create, name)
	defer logEvents()

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
		addEvent(logger.Failure, err)
		return nil
	}
	if res.Server == nil {
		addEvent(logger.Failure, "Empty Response")
		return nil
	}
	if err := ctx.Client.Action.WaitFor(ctx.Context, res.Action); err != nil {
		addEvent(logger.Failure, err)
		return nil
	}
	if err := ctx.Client.Action.WaitFor(ctx.Context, res.NextActions...); err != nil {
		addEvent(logger.Failure, err)
		return nil
	}

	server, err = Get(ctx)
	if err != nil {
		addEvent(logger.Failure, err)
		return nil
	}

	addEvent(logger.Success)
	return server
}
