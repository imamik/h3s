package sshkey

import (
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"log"
)

func Delete(ctx clustercontext.ClusterContext) {
	sshKey := Get(ctx)

	if sshKey == nil {
		return
	}

	log.Println("Deleting ssh key -", sshKey.Name)

	_, err := ctx.Client.SSHKey.Delete(ctx.Context, sshKey)
	if err != nil {
		log.Println("error deleting ssh key:", err)
	}

	log.Println("SSH key deleted -", sshKey.Name)
}
