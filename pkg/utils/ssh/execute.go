package ssh

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/melbahja/goph"
	"hcloud-k3s-cli/pkg/clustercontext"
	"hcloud-k3s-cli/pkg/utils/file"
	"log"
)

func Execute(
	ctx clustercontext.ClusterContext,
	server *hcloud.Server,
	command string,
) {
	privateKeyPath, _ := file.Normalize(ctx.Config.SSHKeyPaths.PrivateKeyPath)
	auth, err := goph.Key(privateKeyPath, "")
	if err != nil {
		log.Fatal(err)
	}

	ip := server.PublicNet.IPv4.IP.String()

	fmt.Printf("ssh -i %s root@%s %s\n", privateKeyPath, ip, command)
	client, err := goph.NewUnknown("root", ip, auth)
	if err != nil {
		log.Fatal(err)
	}
	if client == nil {
		log.Fatal("client is nil")
	}

	defer func(client *goph.Client) {
		err := client.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(client)

	// Execute your command.
	out, err := client.Run(command)

	if err != nil {
		log.Fatal(err)
	}

	// Get your output as []byte.
	fmt.Println(string(out))

}
