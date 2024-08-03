package sshkey

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/clustercontext"
	"h3s/internal/utils/logger"
	"h3s/internal/utils/ssh"
)

func create(ctx clustercontext.ClusterContext) *hcloud.SSHKey {
	publicKey, err := ssh.ReadPublicKeyFromFile(ctx)
	if err != nil {
		logger.LogError(err)
	}

	sshKeyName := getName(ctx)
	addEvent, logEvents := logger.NewEventLogger(logger.SSHKey, logger.Create, sshKeyName)
	defer logEvents()

	sshKey, _, err := ctx.Client.SSHKey.Create(ctx.Context, hcloud.SSHKeyCreateOpts{
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

func Create(ctx clustercontext.ClusterContext) *hcloud.SSHKey {
	sshKey := Get(ctx)
	if sshKey == nil {
		sshKey = create(ctx)
	}
	return sshKey
}
