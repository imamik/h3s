package gateway

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/hetzner/microos/image"
	"h3s/internal/hetzner/network"
	"h3s/internal/hetzner/sshkey"
	"h3s/internal/utils/ip"
	"h3s/internal/utils/logger"
	"h3s/internal/utils/ssh"
)

func Create(ctx *cluster.Cluster) *hcloud.Server {
	net := network.Get(ctx)
	sshKey := sshkey.Get(ctx)

	img, _ := image.Get(ctx, hcloud.ArchitectureARM)
	gateway := createServer(ctx, sshKey, net, img)
	configureGateway(ctx, gateway)
	setupGatewayRoute(ctx, net, gateway)

	return gateway
}

func configureGateway(ctx *cluster.Cluster, gateway *hcloud.Server) {
	ssh.ExecuteWithSsh(ctx, gateway, `
echo 1 > /proc/sys/net/ipv4/ip_forward
iptables -t nat -A POSTROUTING -s '10.0.0.0/16' -o eth0 -j MASQUERADE
	`)
}

func setupGatewayRoute(ctx *cluster.Cluster, net *hcloud.Network, gateway *hcloud.Server) {
	ctx.CloudClient.Network.AddRoute(ctx.Context, net, hcloud.NetworkAddRouteOpts{
		Route: hcloud.NetworkRoute{
			Destination: ip.GetIpRange("0.0.0.0/0"),
			Gateway:     gateway.PrivateNet[0].IP,
		},
	})
}

func createServer(
	ctx *cluster.Cluster,
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

	res, _, err := ctx.CloudClient.Server.Create(ctx.Context, hcloud.ServerCreateOpts{
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
	if err := ctx.CloudClient.Action.WaitFor(ctx.Context, res.Action); err != nil {
		addEvent(logger.Failure, err)
		return nil
	}
	if err := ctx.CloudClient.Action.WaitFor(ctx.Context, res.NextActions...); err != nil {
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
