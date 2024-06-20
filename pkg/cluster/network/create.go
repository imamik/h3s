package network

import (
	"context"
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/utils"
	"hcloud-k3s-cli/pkg/config"
	"log"
	"net"
)

func Create(ctx context.Context, client *hcloud.Client, conf config.Config) (*hcloud.Network, error) {
	_, network, err := net.ParseCIDR("10.0.0.0/16")
	if err != nil {
		log.Fatalf("invalid IP range: %s", err)
	}

	networkResp, _, err := client.Network.Create(ctx, hcloud.NetworkCreateOpts{
		Name:    getName(conf),
		IPRange: network,
		Labels:  utils.GetLabels(conf),
	})
	if err != nil {
		log.Fatalf("error creating network: %s", err)
	}

	fmt.Printf("Created network: %s\n", networkResp.Name)
	return networkResp, nil
}
