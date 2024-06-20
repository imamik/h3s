package server

import (
	"context"
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/config"
	"log"
)

func Create(ctx context.Context, client *hcloud.Client, conf config.Config, networkResp *hcloud.Network, placementGroupResp hcloud.PlacementGroupCreateResult) error {
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
