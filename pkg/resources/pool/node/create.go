package node

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/clustercontext"
	"hcloud-k3s-cli/pkg/config"
	"hcloud-k3s-cli/pkg/utils/logger"
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
	logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Initialized)

	image := &hcloud.Image{Name: "ubuntu-24.04"}
	serverType := &hcloud.ServerType{Name: string(pool.Instance)}
	location := &hcloud.Location{Name: string(pool.Location)}
	publicNet := &hcloud.ServerCreatePublicNet{
		EnableIPv4: pool.EnableIPv4,
		EnableIPv6: pool.EnableIPv6,
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
			"pool":             pool.Name,
			"node":             strconv.Itoa(i),
			"is_control_plane": strconv.FormatBool(isControlPlane),
			"is_worker":        strconv.FormatBool(isWorker),
		}),
	})
	if err != nil {
		logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Failure, err)
	}
	if server.Server == nil {
		logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Failure, "Empty Response")
	}

	logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Success)
	return server.Server
}
