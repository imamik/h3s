package ssh

import (
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/file"
	"hcloud-k3s-cli/internal/utils/logger"
)

func readCertFromFile(path string) (string, error) {
	cert, err := file.Load(path)
	if err != nil {
		logger.LogError(err)
		return "", err
	}
	return string(cert), nil
}

func ReadPublicKeyFromFile(ctx clustercontext.ClusterContext) (string, error) {
	return readCertFromFile(ctx.Config.SSHKeyPaths.PublicKeyPath)
}

func ReadPrivateKeyFromFile(ctx clustercontext.ClusterContext) (string, error) {
	return readCertFromFile(ctx.Config.SSHKeyPaths.PrivateKeyPath)
}