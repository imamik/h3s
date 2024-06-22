package server

import (
	"context"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/utils"
	"hcloud-k3s-cli/pkg/config"
	"log"
	"strconv"
)

func Create(
	name string,
	isControlPlane bool,
	isWorker bool,
	pool config.NodePool,
	network *hcloud.Network,
	placementGroup hcloud.PlacementGroupCreateResult,
	conf config.Config,
	client *hcloud.Client,
	ctx context.Context,
) hcloud.ServerCreateResult {
	name = utils.GetName(name, conf)
	image := &hcloud.Image{Name: "ubuntu-24.04"}
	serverType := &hcloud.ServerType{Name: string(pool.Instance)}
	location := &hcloud.Location{Name: string(pool.Location)}
	publicNet := &hcloud.ServerCreatePublicNet{
		EnableIPv4: false,
		EnableIPv6: false,
	}
	networks := []*hcloud.Network{network}
	sshKeys := []*hcloud.SSHKey{}

	server, _, err := client.Server.Create(ctx, hcloud.ServerCreateOpts{
		Name:           name,
		ServerType:     serverType,
		Image:          image,
		Location:       location,
		Networks:       networks,
		PlacementGroup: placementGroup.PlacementGroup,
		SSHKeys:        sshKeys,
		PublicNet:      publicNet,
		Labels: utils.GetLabels(conf, map[string]string{
			"pool":           pool.Name,
			"isControlPlane": strconv.FormatBool(isControlPlane),
			"isWorker":       strconv.FormatBool(isWorker),
		}),
	})
	if err != nil {
		log.Fatalf("error creating server %s: %s", name, err)
	}
	return server
}
