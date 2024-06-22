package server

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/config"
	"log"
	"strconv"
)

func Create(
	pool config.NodePool,
	i int,
	isControlPlane bool,
	isWorker bool,
	placementGroup hcloud.PlacementGroupCreateResult,
	ctx clustercontext.ClusterContext,
) hcloud.ServerCreateResult {
	name := getName(pool, i, ctx)
	image := &hcloud.Image{Name: "ubuntu-24.04"}
	serverType := &hcloud.ServerType{Name: string(pool.Instance)}
	location := &hcloud.Location{Name: string(pool.Location)}
	publicNet := &hcloud.ServerCreatePublicNet{
		EnableIPv4: false,
		EnableIPv6: false,
	}
	networks := []*hcloud.Network{ctx.Network}
	sshKeys := []*hcloud.SSHKey{}

	server, _, err := ctx.Client.Server.Create(ctx.Context, hcloud.ServerCreateOpts{
		Name:           name,
		ServerType:     serverType,
		Image:          image,
		Location:       location,
		Networks:       networks,
		PlacementGroup: placementGroup.PlacementGroup,
		SSHKeys:        sshKeys,
		PublicNet:      publicNet,
		Labels: ctx.GetLabels(map[string]string{
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
