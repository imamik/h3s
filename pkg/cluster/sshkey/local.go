package sshkey

import (
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/utils/file"
	"hcloud-k3s-cli/pkg/utils/logger"
)

func readLocalPublicKeyFromFile(ctx clustercontext.ClusterContext) (string, error) {
	publicKey, err := file.Load(ctx.Config.SSHKeyPaths.PublicKeyPath)
	if err != nil {
		logger.LogError(err)
		return "", err
	}
	return string(publicKey), nil
}
