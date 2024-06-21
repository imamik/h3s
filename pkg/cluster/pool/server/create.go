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
	network *hcloud.Network,
	placementGroup hcloud.PlacementGroupCreateResult,
	conf config.Config,
	client *hcloud.Client,
	ctx context.Context,
) hcloud.ServerCreateResult {
	image := &hcloud.Image{Name: "ubuntu-20.04"}
	serverType := &hcloud.ServerType{Name: "cax11"}

	server, _, err := client.Server.Create(ctx, hcloud.ServerCreateOpts{
		Name:           utils.GetName(name, conf),
		ServerType:     serverType,
		Image:          image,
		SSHKeys:        []*hcloud.SSHKey{},
		Networks:       []*hcloud.Network{network},
		PlacementGroup: placementGroup.PlacementGroup,
		Labels: utils.GetLabels(conf, map[string]string{
			"isControlPlane": strconv.FormatBool(isControlPlane),
			"isWorker":       strconv.FormatBool(isWorker),
		}),
	})
	if err != nil {
		log.Fatalf("error creating server %s: %s", name, err)
	}
	return server
}
