package sshkey

import (
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/utils/file"
	"log"
)

func readLocalPublicKeyFromFile(ctx clustercontext.ClusterContext) (string, error) {
	publicKey, err := file.Load(ctx.Config.SSHKeyPaths.PublicKeyPath)
	if err != nil {
		log.Println("error reading public key from file:", err)
		return "", err
	}
	return string(publicKey), nil
}
