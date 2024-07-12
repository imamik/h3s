package ssh

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"golang.org/x/crypto/ssh"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/ip"
	"log"
	"time"
)

func ExecuteWithSsh(
	ctx clustercontext.ClusterContext,
	remote *hcloud.Server,
	command string,
) (string, error) {
	remoteIp := ip.FirstAvailable(remote)
	removeKnownHostsEntry(remoteIp)

	// SSH client configuration
	sshConfig, err := ConfigSsh(ctx)
	if err != nil {
		return "", fmt.Errorf("unable to create SSH config: %w", err)
	}

	// Connect to the remote server
	client, err := dialWithRetries(remoteIp, sshConfig, 5*time.Second, 5)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}

	res, err := Run(client, command)
	if err != nil {
		return "", err
	}

	return res, nil
}

func ExecuteViaProxy(
	ctx clustercontext.ClusterContext,
	gateway *hcloud.Server,
	remote *hcloud.Server,
	command string,
) (string, error) {
	if gateway == nil {
		return ExecuteWithSsh(ctx, remote, command)
	}

	proxyIp := ip.FirstAvailable(gateway)
	remoteIp := ip.FirstAvailable(remote)
	removeKnownHostsEntry(proxyIp)

	sshConfig, err := ConfigSsh(ctx)
	if err != nil {
		return "", fmt.Errorf("unable to create SSH config: %w", err)
	}

	proxyConn, err := dialWithRetries(proxyIp, sshConfig, 5*time.Second, 5)
	if err != nil {
		return "", fmt.Errorf("unable to connect to gateway: %w", err)
	}
	defer func(proxyConn *ssh.Client) {
		err := proxyConn.Close()
		if err != nil {
			fmt.Printf("unable to close gateway connection: %v\n", err)
		}
	}(proxyConn)

	tunnel, err := proxyConn.Dial("tcp", remoteIp+":22")
	if err != nil {
		return "", fmt.Errorf("unable to create tunnel to remote: %w", err)
	}

	clientConn, chans, reqs, err := ssh.NewClientConn(tunnel, remoteIp+":22", sshConfig)
	if err != nil {
		return "", fmt.Errorf("unable to establish connection to remote: %w", err)
	}
	client := ssh.NewClient(clientConn, chans, reqs)

	res, err := Run(client, command)
	if err != nil {
		return "", err
	}

	return res, nil

}
