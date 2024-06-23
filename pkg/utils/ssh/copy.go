package ssh

import (
	"fmt"
	"github.com/melbahja/goph"
	"hcloud-k3s-cli/pkg/clustercontext"
	"strings"
)

func echoTo(key string, path string) string {
	return fmt.Sprintf("echo '%s' > %s", key, path)
}

func adjustPermissions(path string) string {
	return fmt.Sprintf("chmod 600 %s", path)
}

func uploadCertsToServer(ctx clustercontext.ClusterContext, client *goph.Client) {
	fmt.Println("Uploading SSH keys to server")

	privateKeyPath := ctx.Config.SSHKeyPaths.PrivateKeyPath
	publicKeyPath := ctx.Config.SSHKeyPaths.PublicKeyPath
	privateKey, _ := ReadPrivateKeyFromFile(ctx)
	publicKey, _ := ReadPublicKeyFromFile(ctx)

	cmdArr := []string{
		echoTo(privateKey, privateKeyPath),
		adjustPermissions(privateKeyPath),
		echoTo(publicKey, publicKeyPath),
		adjustPermissions(publicKeyPath),
	}
	cmd := strings.Join(cmdArr, " && ")

	run(client, cmd)
}
