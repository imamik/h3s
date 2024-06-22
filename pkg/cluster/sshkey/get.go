package sshkey

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"log"
)

func Get(ctx clustercontext.ClusterContext) *hcloud.SSHKey {
	sshKeyName := getName(ctx)
	log.Println("Getting ssh key - " + sshKeyName)

	sshKey, _, err := ctx.Client.SSHKey.GetByName(ctx.Context, sshKeyName)
	if err != nil {
		log.Println("error getting ssh key:", err)
	}
	if sshKey == nil {
		log.Println("ssh key not found:", sshKeyName)
		return nil
	}

	log.Println("SSH key found - " + sshKey.Name)
	return sshKey
}
