package loadbalancer

import (
	"context"
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/config"
	"log"
)

func Create(ctx context.Context, client *hcloud.Client, conf config.Config, network *hcloud.Network) error {
	loadBalancerResp, _, err := client.LoadBalancer.Create(ctx, hcloud.LoadBalancerCreateOpts{
		Name:             "lb-kubernetes",
		LoadBalancerType: &hcloud.LoadBalancerType{Name: "lb11"},
		Location:         &hcloud.Location{Name: "nbg1"},
		Network:          network,
	})
	if err != nil {
		log.Fatalf("error creating load balancer: %s", err)
	}
	fmt.Printf("Created load balancer: %s\n", loadBalancerResp.LoadBalancer.Name)
	return nil
}
