package sshkey

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Get(ctx clustercontext.ClusterContext) *hcloud.SSHKey {
	sshKeyName := getName(ctx)

	addEvent, logEvents := logger.NewEventLogger(logger.SSHKey, logger.Get, sshKeyName)
	defer logEvents()

	sshKey, _, err := ctx.Client.SSHKey.GetByName(ctx.Context, sshKeyName)
	if err != nil || sshKey == nil {
		addEvent(logger.Failure, err)
		return nil
	}

	addEvent(logger.Success)
	return sshKey
}
