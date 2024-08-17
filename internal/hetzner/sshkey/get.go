package sshkey

import (
	"errors"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
)

func Get(ctx *cluster.Cluster) (*hcloud.SSHKey, error) {
	sshKeyName := getName(ctx)

	l := logger.New(nil, logger.SSHKey, logger.Get, sshKeyName)
	defer l.LogEvents()

	sshKey, _, err := ctx.CloudClient.SSHKey.GetByName(ctx.Context, sshKeyName)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}
	if sshKey == nil {
		err = errors.New("ssh key is nil")
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	l.AddEvent(logger.Success)
	return sshKey, nil
}
