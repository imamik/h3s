package gateway

import (
	"errors"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/hetzner/microos/image"
	"h3s/internal/hetzner/network"
	"h3s/internal/hetzner/sshkey"
	"h3s/internal/utils/ip"
	"h3s/internal/utils/logger"
	"h3s/internal/utils/ssh"
	"net"
)

func Create(ctx *cluster.Cluster) (*hcloud.Server, error) {
	l := logger.New(nil, logger.Server, logger.Create, "gateway")
	defer l.LogEvents()

	n, err := network.Get(ctx)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	// Get ssh key
	sshKey, err := sshkey.Get(ctx)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	// Get image for architecture
	img, err := image.Get(ctx, hcloud.ArchitectureARM)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}
	gateway, err := createServer(ctx, sshKey, n, img)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	// Configure gateway
	_, err = configureGateway(ctx, gateway)
	if err != nil {
		return nil, err
	}

	// Setup route for gateway
	err = setupGatewayRoute(ctx, n, gateway)
	if err != nil {
		return nil, err
	}

	return gateway, nil
}

func configureGateway(ctx *cluster.Cluster, gateway *hcloud.Server) (string, error) {
	return ssh.ExecuteWithSsh(ctx.Config.SSHKeyPaths.PrivateKeyPath, gateway, `
echo 1 > /proc/sys/net/ipv4/ip_forward
iptables -t nat -A POSTROUTING -s '10.0.0.0/16' -o eth0 -j MASQUERADE
	`)
}

func setupGatewayRoute(ctx *cluster.Cluster, network *hcloud.Network, gateway *hcloud.Server) error {
	_, ipRange, err := net.ParseCIDR("0.0.0.0/0")
	if err != nil {
		return err
	}
	_, _, err = ctx.CloudClient.Network.AddRoute(ctx.Context, network, hcloud.NetworkAddRouteOpts{
		Route: hcloud.NetworkRoute{
			Destination: ipRange,
			Gateway:     ip.Private(gateway),
		},
	})
	return err
}

func createServer(
	ctx *cluster.Cluster,
	sshKey *hcloud.SSHKey,
	network *hcloud.Network,
	image *hcloud.Image,
) (*hcloud.Server, error) {
	name := getName(ctx)
	l := logger.New(nil, logger.Server, logger.Create, name)
	defer l.LogEvents()

	server, err := Get(ctx)
	if server != nil && err == nil {
		l.AddEvent(logger.Success, "gateway already exists")
		return server, nil
	}

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

	server, err = Get(ctx)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	l.AddEvent(logger.Success)
	return server, nil
}
