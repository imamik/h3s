package sshkey

import (
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/utils/logger"
)

func Delete(ctx clustercontext.ClusterContext) {
	sshKey := Get(ctx)

	if sshKey == nil {
		return
	}

	logger.LogResourceEvent(logger.SSHKey, logger.Delete, sshKey.Name, logger.Initialized)

	_, err := ctx.Client.SSHKey.Delete(ctx.Context, sshKey)
	if err != nil {
		logger.LogResourceEvent(logger.SSHKey, logger.Delete, sshKey.Name, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.SSHKey, logger.Delete, sshKey.Name, logger.Success)
}
