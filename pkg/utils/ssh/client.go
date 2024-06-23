package ssh

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/melbahja/goph"
	"hcloud-k3s-cli/pkg/clustercontext"
	"hcloud-k3s-cli/pkg/utils/file"
	"log"
)

func Client(
	ctx clustercontext.ClusterContext,
	gate *hcloud.Server,
) *goph.Client {
	gateWayIp := gate.PublicNet.IPv4.IP.String()
	privateKeyPath := ctx.Config.SSHKeyPaths.PrivateKeyPath
	client := getClient(gateWayIp, privateKeyPath, "")
	uploadCertsToServer(ctx, client)
	return client
}

func getClient(
	ip string,
	privateKeyPath string,
	passphrase string,
) *goph.Client {
	keyPath, _ := file.Normalize(privateKeyPath)
	auth, err := goph.Key(keyPath, passphrase)
	if err != nil {
		log.Fatal(err)
	}

	client, err := goph.NewUnknown("root", ip, auth)
	if err != nil {
		log.Fatal(err)
	}
	if client == nil {
		log.Fatal("client is nil")
	}

	return client
}
