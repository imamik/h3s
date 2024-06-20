package cluster

import (
	"context"
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/config"
	"log"
	"net"
	"os"
)

func getName(conf config.Config, name string) string {
	return fmt.Sprintf("%s-k3s-cluster-%s", conf.Name, name)
}

func createNetwork(ctx context.Context, client *hcloud.Client, conf config.Config) (*hcloud.Network, error) {
	ipRange := "10.0.0.0/16"
	_, network, err := net.ParseCIDR(ipRange)
	if err != nil {
		log.Fatalf("invalid IP range: %s", err)
	}
	// Prefix the network name with the project name to avoid conflicts
	networkResp, _, err := client.Network.Create(ctx, hcloud.NetworkCreateOpts{
		Name:    getName(conf, "network"),
		IPRange: network,
		Labels: map[string]string{
			"project": conf.Name,
		},
	})
	if err != nil {
		log.Fatalf("error creating network: %s", err)
	}
	fmt.Printf("Created network: %s\n", networkResp.Name)
	return networkResp, nil
}

func createPlacementGroup(ctx context.Context, client *hcloud.Client, conf config.Config) (hcloud.PlacementGroupCreateResult, error) {
	placementGroupResp, _, err := client.PlacementGroup.Create(ctx, hcloud.PlacementGroupCreateOpts{
		Name: getName(conf, "placement-group"),
		Type: hcloud.PlacementGroupTypeSpread,
		Labels: map[string]string{
			"project": conf.Name,
		},
	})
	if err != nil {
		log.Fatalf("error creating placement group: %s", err)
	}
	fmt.Printf("Created placement group: %s\n", placementGroupResp.PlacementGroup.Name)
	return placementGroupResp, nil
}

func createServers(ctx context.Context, client *hcloud.Client, conf config.Config, networkResp *hcloud.Network, placementGroupResp hcloud.PlacementGroupCreateResult) error {
	serverNames := []string{"k3s-node-1", "k3s-node-2", "k3s-node-3"}
	for _, name := range serverNames {
		_, _, err := client.Server.Create(ctx, hcloud.ServerCreateOpts{
			Name:           name,
			ServerType:     &hcloud.ServerType{Name: "cx11"},
			Image:          &hcloud.Image{Name: "ubuntu-20.04"},
			SSHKeys:        []*hcloud.SSHKey{},
			Networks:       []*hcloud.Network{networkResp},
			PlacementGroup: placementGroupResp.PlacementGroup,
			Labels: map[string]string{
				"role": "kubernetes",
			},
		})
		if err != nil {
			log.Fatalf("error creating server %s: %s", name, err)
		}
		fmt.Printf("Created server: %s\n", name)
	}
	return nil
}

func createFirewall(ctx context.Context, client *hcloud.Client, conf config.Config, network *net.IPNet) error {
	firewallResp, _, err := client.Firewall.Create(ctx, hcloud.FirewallCreateOpts{
		Name: "kubernetes-firewall",
		Rules: []hcloud.FirewallRule{
			{
				Direction:   hcloud.FirewallRuleDirectionIn,
				Protocol:    hcloud.FirewallRuleProtocolTCP,
				Port:        hcloud.Ptr("22"),
				SourceIPs:   []net.IPNet{*network},
				Description: hcloud.Ptr("SSH"),
			},
			{
				Direction:   hcloud.FirewallRuleDirectionIn,
				Protocol:    hcloud.FirewallRuleProtocolTCP,
				Port:        hcloud.Ptr("6443"),
				SourceIPs:   []net.IPNet{*network},
				Description: hcloud.Ptr("K3S API Server"),
			},
			{
				Direction:   hcloud.FirewallRuleDirectionIn,
				Protocol:    hcloud.FirewallRuleProtocolICMP,
				SourceIPs:   []net.IPNet{*network},
				Description: hcloud.Ptr("ICMP"),
			},
		},
	})
	if err != nil {
		log.Fatalf("error creating firewall: %s", err)
	}
	fmt.Printf("Created firewall: %s\n", firewallResp.Firewall.Name)
	return nil
}

func createLoadBalancer(ctx context.Context, client *hcloud.Client, conf config.Config, network *hcloud.Network) error {
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

func CreateCluster(conf config.Config) {
	ctx := context.Background()
	client := hcloud.NewClient(hcloud.WithToken(os.Getenv("HCLOUD_TOKEN")))

	// Step 1: Create Private Network
	networkResp, networkErr := createNetwork(ctx, client, conf)
	if networkErr != nil {
		log.Fatalf("error creating network: %s", networkErr)
	}
	fmt.Println(networkResp)

	// Step 2: Create Placement Group
	placementGroupResp, placementGroupErr := createPlacementGroup(ctx, client, conf)
	if placementGroupErr != nil {
		log.Fatalf("error creating placement group: %s", placementGroupErr)
	}
	fmt.Println(placementGroupResp)

	// Step 3: Create Servers
	//err = createServers(ctx, client, conf, networkResp, placementGroupResp)

	// Step 4: Create Firewall
	//err = createFirewall(ctx, client, conf, networkResp.IPRange)

	// Step 5: Create Load Balancer
	//err = createLoadBalancer(ctx, client, conf, networkResp)
}
