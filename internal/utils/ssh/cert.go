package ssh

import (
	"h3s/internal/clustercontext"
	"h3s/internal/utils/file"
	"h3s/internal/utils/logger"
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
