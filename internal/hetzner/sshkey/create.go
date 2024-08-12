package sshkey

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/utils/file"
	"h3s/internal/utils/logger"
)

func create(ctx *cluster.Cluster) *hcloud.SSHKey {
	publicKey, err := file.New(ctx.Config.SSHKeyPaths.PublicKeyPath).Load().GetString()
	if err != nil {
		logger.LogError(err)
	}

	sshKeyName := getName(ctx)
	addEvent, logEvents := logger.NewEventLogger(logger.SSHKey, logger.Create, sshKeyName)
	defer logEvents()

	sshKey, _, err := ctx.CloudClient.SSHKey.Create(ctx.Context, hcloud.SSHKeyCreateOpts{
		Name:      sshKeyName,
		PublicKey: publicKey,
		Labels:    ctx.GetLabels(),
	})
	if err != nil || sshKey == nil {
		addEvent(logger.Failure, err)
		return nil
	}

	addEvent(logger.Success)
	return sshKey
}

func Create(ctx *cluster.Cluster) *hcloud.SSHKey {
	sshKey := Get(ctx)
	if sshKey == nil {
		sshKey = create(ctx)
	}
	return sshKey
}
