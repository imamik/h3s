package sshkey

import (
	"errors"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/utils/file"
	"h3s/internal/utils/logger"
)

func create(ctx *cluster.Cluster) (*hcloud.SSHKey, error) {
	publicKey, err := file.New(ctx.Config.SSHKeyPaths.PublicKeyPath).Load().GetString()

	l := logger.New(nil, logger.SSHKey, logger.Create, getName(ctx))
	defer l.LogEvents()
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	sshKeyName := getName(ctx)

	sshKey, _, err := ctx.CloudClient.SSHKey.Create(ctx.Context, hcloud.SSHKeyCreateOpts{
		Name:      sshKeyName,
		PublicKey: publicKey,
		Labels:    ctx.GetLabels(),
	})
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}
	if sshKey == nil {
		err = errors.New("sshKey is nil")
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	l.AddEvent(logger.Success)
	return sshKey, nil
}

func Create(ctx *cluster.Cluster) (*hcloud.SSHKey, error) {
	sshKey, err := Get(ctx)
	if sshKey != nil && err == nil {
		return sshKey, nil
	}
	return create(ctx)
}