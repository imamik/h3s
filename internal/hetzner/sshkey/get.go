package sshkey

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
)

func Get(ctx *cluster.Cluster) *hcloud.SSHKey {
	sshKeyName := getName(ctx)

	addEvent, logEvents := logger.NewEventLogger(logger.SSHKey, logger.Get, sshKeyName)
	defer logEvents()

	sshKey, _, err := ctx.CloudClient.SSHKey.GetByName(ctx.Context, sshKeyName)
	if err != nil || sshKey == nil {
		addEvent(logger.Failure, err)
		return nil
	}

	addEvent(logger.Success)
	return sshKey
}
