package node

import (
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
	addEvent, logEvents := logger.NewEventLogger(logger.Server, logger.Create, name)
	defer logEvents()

	serverType := hcloud.ServerType{Name: string(pool.Instance)}
	location := hcloud.Location{Name: string(pool.Location)}
	publicNet := hcloud.ServerCreatePublicNet{
		EnableIPv4: ctx.Config.PublicIps,
		EnableIPv6: ctx.Config.PublicIps,
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
		addEvent(logger.Failure, err)
		return nil
	}
	if res.Server == nil {
		addEvent(logger.Failure, "Empty Response")
		return nil
	}
	if err := ctx.Client.Action.WaitFor(ctx.Context, res.Action); err != nil {
		addEvent(logger.Failure, err)
		return nil
	}
	if err := ctx.Client.Action.WaitFor(ctx.Context, res.NextActions...); err != nil {
		addEvent(logger.Failure, err)
		return nil
	}

	addEvent(logger.Success)
	return res.Server
}
