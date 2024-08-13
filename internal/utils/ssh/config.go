package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"h3s/internal/cluster"
	"h3s/internal/utils/file"
)

const (
	user = "root"
)

func ConfigSsh(ctx *cluster.Cluster) (*ssh.ClientConfig, error) {
	// Load private key
	key, err := file.New(ctx.Config.SSHKeyPaths.PrivateKeyPath).Load().GetBytes()

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %w", err)
	}

	// SSH client configuration
	return &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, nil
}
