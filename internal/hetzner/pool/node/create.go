// Package node contains the functionality for managing Hetzner cloud servers
package node

import (
	"errors"
	"h3s/internal/cluster"
	"h3s/internal/config"
	"h3s/internal/hetzner/pool/node/userdata"
	"h3s/internal/utils/logger"
	"strconv"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Create creates a Hetzner cloud server
func Create(
	ctx *cluster.Cluster,
	sshKey *hcloud.SSHKey,
	network *hcloud.Network,
	image *hcloud.Image,
	placementGroup *hcloud.PlacementGroup,
	pool config.NodePool,
	i int,
	isControlPlane bool,
	isWorker bool,
) (*hcloud.Server, error) {
	server, err := Get(ctx, pool, i)
	if server != nil && err == nil {
		return server, nil
	}
	return create(ctx, sshKey, network, image, placementGroup, pool, i, isControlPlane, isWorker)
}

func create(
	ctx *cluster.Cluster,
	sshKey *hcloud.SSHKey,
	network *hcloud.Network,
	image *hcloud.Image,
	placementGroup *hcloud.PlacementGroup,
	pool config.NodePool,
	i int,
	isControlPlane bool,
	isWorker bool,
) (*hcloud.Server, error) {
	name := getName(ctx, pool, i)

	l := logger.New(nil, logger.Server, logger.Create, name)
	defer l.LogEvents()

	serverType := hcloud.ServerType{Name: string(pool.Instance)}
	location := hcloud.Location{Name: string(pool.Location)}
	publicNet := hcloud.ServerCreatePublicNet{
		EnableIPv4: false,
		EnableIPv6: false,
	}
	networks := []*hcloud.Network{network}
	sshKeys := []*hcloud.SSHKey{sshKey}

	usrData, err := userdata.GenerateCloudInitConfig(userdata.CloudInitConfig{
		Hostname:        name,
		SSHPort:         22,
		SSHMaxAuthTries: 5,
		SSHAuthorizedKeys: []string{
			sshKey.PublicKey,
		},
	})
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	res, _, err := ctx.CloudClient.Server.Create(ctx.Context, hcloud.ServerCreateOpts{
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
		l.AddEvent(logger.Failure, err)
		return nil, err
	}
	if res.Server == nil {
		err = errors.New("server is nil")
		l.AddEvent(logger.Failure, err)
		return nil, err
	}
	if err := ctx.CloudClient.Action.WaitFor(ctx.Context, res.Action); err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}
	if err := ctx.CloudClient.Action.WaitFor(ctx.Context, res.NextActions...); err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	l.AddEvent(logger.Success)
	return res.Server, nil
}
