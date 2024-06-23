package ssh

import (
	"github.com/melbahja/goph"
	"hcloud-k3s-cli/pkg/utils/file"
	"log"
)

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
