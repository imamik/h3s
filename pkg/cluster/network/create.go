package network

import (
	"context"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/utils"
	"hcloud-k3s-cli/pkg/config"
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

func Create(ctx context.Context, client *hcloud.Client, conf config.Config) {
	network, _, err := client.Network.Create(ctx, hcloud.NetworkCreateOpts{
		Name:    getName(conf),
		IPRange: getIpRange("10.0.0.0/16"),
		Labels:  utils.GetLabels(conf),
	})
	if err != nil {
		log.Fatalf("error creating network: %s", err)
	}

	subnet := hcloud.NetworkSubnet{
		Type:        hcloud.NetworkSubnetTypeServer,
		IPRange:     getIpRange("10.0.0.0/16"),
		NetworkZone: conf.NetworkZone,
	}

	_, _, err = client.Network.AddSubnet(ctx, network, hcloud.NetworkAddSubnetOpts{
		Subnet: subnet,
	})
	if err != nil {
		log.Fatalf("error adding subnet to network: %s", err)
	}
}
