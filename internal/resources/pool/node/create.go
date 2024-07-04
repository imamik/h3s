package node

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/config"
	"hcloud-k3s-cli/internal/resources/pool/node/userdata"
	"hcloud-k3s-cli/internal/utils/logger"
	"strconv"
)

func Create(
	ctx clustercontext.ClusterContext,
	sshKey *hcloud.SSHKey,
	network *hcloud.Network,
	image *hcloud.Image,
	placementGroup *hcloud.PlacementGroup,
	pool config.NodePool,
	i int,
	isControlPlane bool,
	isWorker bool,
) *hcloud.Server {
	server := Get(ctx, pool, i)
	if server == nil {
		server = create(ctx, sshKey, network, image, placementGroup, pool, i, isControlPlane, isWorker)
	}
	return server
}

func create(
	ctx clustercontext.ClusterContext,
	sshKey *hcloud.SSHKey,
	network *hcloud.Network,
	image *hcloud.Image,
	placementGroup *hcloud.PlacementGroup,
	pool config.NodePool,
	i int,
	isControlPlane bool,
	isWorker bool,
) *hcloud.Server {
	name := getName(ctx, pool, i)
	logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Initialized)

	serverType := hcloud.ServerType{Name: string(pool.Instance)}
	location := hcloud.Location{Name: string(pool.Location)}
	publicNet := hcloud.ServerCreatePublicNet{
		EnableIPv4: pool.EnableIPv4,
		EnableIPv6: pool.EnableIPv6,
	}
	networks := []*hcloud.Network{network}
	sshKeys := []*hcloud.SSHKey{sshKey}

	usrData := userdata.GenerateCloudInitConfig(userdata.CloudInitConfig{
		Hostname:        name,
		SSHPort:         22,
		SSHMaxAuthTries: 5,
		SSHAuthorizedKeys: []string{
			sshKey.PublicKey,
		},
	})

	fmt.Printf("\n\n====================================\nUser Data:\n\n%s\n\n====================================\n\n", usrData)

	res, _, err := ctx.Client.Server.Create(ctx.Context, hcloud.ServerCreateOpts{
		Name:           name,
		ServerType:     &serverType,
		Image:          image,
		Location:       &location,
		Networks:       networks,
		PlacementGroup: placementGroup,
		SSHKeys:        sshKeys,
		PublicNet:      &publicNet,
		UserData:       usrData,
		Labels: ctx.GetLabels(map[string]string{
			"pool":                 pool.Name,
			"node":                 strconv.Itoa(i),
			string(IsControlPlane): strconv.FormatBool(isControlPlane),
			string(IsWorker):       strconv.FormatBool(isWorker),
		}),
	})
	if err != nil {
		logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Failure, err)
	}
	if res.Server == nil {
		logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Failure, "Empty Response")
	}
	if err := ctx.Client.Action.WaitFor(ctx.Context, res.Action); err != nil {
		logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Failure, err)
	}
	if err := ctx.Client.Action.WaitFor(ctx.Context, res.NextActions...); err != nil {
		logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.Server, logger.Create, name, logger.Success)
	return res.Server
}
