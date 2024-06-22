package sshkey

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/resources/clustercontext"
	"hcloud-k3s-cli/pkg/utils/logger"
)

func create(ctx clustercontext.ClusterContext) *hcloud.SSHKey {
	publicKey, err := readLocalPublicKeyFromFile(ctx)
	if err != nil {
		logger.LogError(err)
	}

	sshKeyName := getName(ctx)
	logger.LogResourceEvent(logger.SSHKey, logger.Create, sshKeyName, logger.Initialized)

	sshKey, _, err := ctx.Client.SSHKey.Create(ctx.Context, hcloud.SSHKeyCreateOpts{
		Name:      sshKeyName,
		PublicKey: publicKey,
		Labels:    ctx.GetLabels(),
	})
	if err != nil || sshKey == nil {
		logger.LogResourceEvent(logger.SSHKey, logger.Create, sshKeyName, logger.Failure, err)
	}

	logger.LogResourceEvent(logger.SSHKey, logger.Create, sshKeyName, logger.Success)
	return sshKey
}

func Create(ctx clustercontext.ClusterContext) *hcloud.SSHKey {
	sshKey := Get(ctx)
	if sshKey == nil {
		sshKey = create(ctx)
	}
	return sshKey
}
