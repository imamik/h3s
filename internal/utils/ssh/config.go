package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"hcloud-k3s-cli/internal/clustercontext"
)

const (
	user = "root"
)

func ConfigSsh(ctx clustercontext.ClusterContext) (*ssh.ClientConfig, error) {
	// Load private key
	key, err := ReadPrivateKeyFromFile(ctx)

	signer, err := ssh.ParsePrivateKey([]byte(key))
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

func ConfigPass(password string) (*ssh.ClientConfig, error) {
	return &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, nil
}
