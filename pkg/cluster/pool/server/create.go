package server

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/config"
	"log"
	"strconv"
)

func Create(
	ctx clustercontext.ClusterContext,
	sshKey *hcloud.SSHKey,
	network *hcloud.Network,
	placementGroup *hcloud.PlacementGroup,
	pool config.NodePool,
	i int,
	isControlPlane bool,
	isWorker bool,
) *hcloud.Server {
	server := Get(ctx, pool, i)
	if server == nil {
		server = create(ctx, sshKey, network, placementGroup, pool, i, isControlPlane, isWorker)
	}
	return server
}

func create(
	ctx clustercontext.ClusterContext,
	sshKey *hcloud.SSHKey,
	network *hcloud.Network,
	placementGroup *hcloud.PlacementGroup,
	pool config.NodePool,
	i int,
	isControlPlane bool,
	isWorker bool,
) *hcloud.Server {
	name := getName(ctx, pool, i)
	log.Println("Creating server - " + name)

	image := &hcloud.Image{Name: "ubuntu-24.04"}
	serverType := &hcloud.ServerType{Name: string(pool.Instance)}
	location := &hcloud.Location{Name: string(pool.Location)}
	publicNet := &hcloud.ServerCreatePublicNet{
		EnableIPv4: false,
		EnableIPv6: false,
	}
	networks := []*hcloud.Network{network}
	sshKeys := []*hcloud.SSHKey{sshKey}

	server, _, err := ctx.Client.Server.Create(ctx.Context, hcloud.ServerCreateOpts{
		Name:           name,
		ServerType:     serverType,
		Image:          image,
		Location:       location,
		Networks:       networks,
		PlacementGroup: placementGroup,
		SSHKeys:        sshKeys,
		PublicNet:      publicNet,
		Labels: ctx.GetLabels(map[string]string{
			"pool":           pool.Name,
			"isControlPlane": strconv.FormatBool(isControlPlane),
			"isWorker":       strconv.FormatBool(isWorker),
		}),
	})
	if err != nil {
		log.Println("error creating server %s: %s", name, err)
	}

	log.Println("Server created - " + name)
	return server.Server
}
