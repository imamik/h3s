package sshkey

import (
	"h3s/internal/cluster"
	"h3s/internal/utils/logger"
)

// Delete deletes a Hetzner cloud SSH key
func Delete(ctx *cluster.Cluster) error {
	l := logger.New(nil, logger.SSHKey, logger.Delete, getName(ctx))
	defer l.LogEvents()

	sshKey, err := Get(ctx)
	if sshKey == nil && err.Error() == "ssh key is nil" {
		l.AddEvent(logger.Success, "no ssh key found to delete")
		return nil
	}

	_, err = ctx.CloudClient.SSHKey.Delete(ctx.Context, sshKey)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	l.AddEvent(logger.Success)
	return nil
}
