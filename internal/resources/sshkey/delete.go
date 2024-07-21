package sshkey

import (
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Delete(ctx clustercontext.ClusterContext) {
	sshKey := Get(ctx)

	if sshKey == nil {
		return
	}

	addEvent, logEvents := logger.NewEventLogger(logger.SSHKey, logger.Delete, sshKey.Name)
	defer logEvents()

	_, err := ctx.Client.SSHKey.Delete(ctx.Context, sshKey)
	if err != nil {
		addEvent(logger.Failure, err)
		return
	}

	addEvent(logger.Success)
}
