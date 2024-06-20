package firewall

import (
	"context"
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/config"
	"log"
	"net"
)

func Create(ctx context.Context, client *hcloud.Client, conf config.Config, network *net.IPNet) error {
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
