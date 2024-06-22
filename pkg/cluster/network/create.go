package network

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"log"
	"net"
)

func Create(ctx clustercontext.ClusterContext) *hcloud.Network {
	network := Get(ctx)
	if network == nil {
		network = create(ctx)
	}
	return network
}

func getIpRange(s string) *net.IPNet {
	_, ipRange, err := net.ParseCIDR(s)
	if err != nil {
		log.Println("invalid IP range: %s", err)
	}
	return ipRange
}

func create(ctx clustercontext.ClusterContext) *hcloud.Network {
	networkName := getName(ctx)
	log.Println("Creating network - ", networkName)

	network, _, err := ctx.Client.Network.Create(ctx.Context, hcloud.NetworkCreateOpts{
		Name:    networkName,
		IPRange: getIpRange("10.0.0.0/16"),
		Labels:  ctx.GetLabels(),
	})
	if err != nil {
		log.Println("error creating network: ", err)
	}

	log.Println("Network created - " + networkName)
	log.Println("Adding subnet to network - " + networkName)

	subnet := hcloud.NetworkSubnet{
		Type:        hcloud.NetworkSubnetTypeServer,
		IPRange:     getIpRange("10.0.0.0/16"),
		NetworkZone: ctx.Config.NetworkZone,
	}

	_, _, err = ctx.Client.Network.AddSubnet(ctx.Context, network, hcloud.NetworkAddSubnetOpts{
		Subnet: subnet,
	})
	if err != nil {
		log.Println("error adding subnet to network: ", err)
	}

	log.Println("Subnet added to network - " + networkName)

	return network
}
