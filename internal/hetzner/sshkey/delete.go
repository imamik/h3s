package sshkey

import (
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
)

func Delete(ctx *cluster.Cluster) {
	sshKey := Get(ctx)

	if sshKey == nil {
		return
	}

	addEvent, logEvents := logger.NewEventLogger(logger.SSHKey, logger.Delete, sshKey.Name)
	defer logEvents()

	_, err := ctx.CloudClient.SSHKey.Delete(ctx.Context, sshKey)
	if err != nil {
		addEvent(logger.Failure, err)
		return
	}

	addEvent(logger.Success)
}
