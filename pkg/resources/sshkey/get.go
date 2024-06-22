package sshkey

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/resources/clustercontext"
	"hcloud-k3s-cli/pkg/utils/logger"
)

func Get(ctx clustercontext.ClusterContext) *hcloud.SSHKey {
	sshKeyName := getName(ctx)
	logger.LogResourceEvent(logger.SSHKey, logger.Get, sshKeyName, logger.Initialized)

	sshKey, _, err := ctx.Client.SSHKey.GetByName(ctx.Context, sshKeyName)
	if err != nil || sshKey == nil {
		logger.LogResourceEvent(logger.SSHKey, logger.Get, sshKeyName, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.SSHKey, logger.Get, sshKeyName, logger.Success)
	return sshKey
}
