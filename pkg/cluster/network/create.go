package network

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"log"
	"net"
)

func getIpRange(s string) *net.IPNet {
	_, ipRange, err := net.ParseCIDR(s)
	if err != nil {
		log.Fatalf("invalid IP range: %s", err)
	}
	return ipRange
}

func Create(ctx clustercontext.ClusterContext) *hcloud.Network {
	network, _, err := ctx.Client.Network.Create(ctx.Context, hcloud.NetworkCreateOpts{
		Name:    getName(ctx),
		IPRange: getIpRange("10.0.0.0/16"),
		Labels:  ctx.GetLabels(),
	})
	if err != nil {
		log.Fatalf("error creating network: %s", err)
	}

	subnet := hcloud.NetworkSubnet{
		Type:        hcloud.NetworkSubnetTypeServer,
		IPRange:     getIpRange("10.0.0.0/16"),
		NetworkZone: ctx.Config.NetworkZone,
	}

	_, _, err = ctx.Client.Network.AddSubnet(ctx.Context, network, hcloud.NetworkAddSubnetOpts{
		Subnet: subnet,
	})
	if err != nil {
		log.Fatalf("error adding subnet to network: %s", err)
	}

	return network
}
