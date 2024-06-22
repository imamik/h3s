package sshkey

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"log"
)

func create(ctx clustercontext.ClusterContext) *hcloud.SSHKey {
	publicKey, err := readLocalPublicKeyFromFile(ctx)
	if err != nil {
		log.Println("error reading public key:", err)
		return nil
	}

	sshKeyName := getName(ctx)
	log.Println("Creating ssh key - " + sshKeyName)

	sshKey, _, err := ctx.Client.SSHKey.Create(ctx.Context, hcloud.SSHKeyCreateOpts{
		Name:      sshKeyName,
		PublicKey: publicKey,
		Labels:    ctx.GetLabels(),
	})
	if err != nil {
		log.Println("error creating ssh key: ", err)
	}
	if sshKey == nil {
		log.Println("ssh key not created")
		return nil
	}

	log.Println("SSH key created - " + sshKey.Name)
	return sshKey
}

func Create(ctx clustercontext.ClusterContext) *hcloud.SSHKey {
	sshKey := Get(ctx)
	if sshKey == nil {
		sshKey = create(ctx)
	}
	return sshKey
}
