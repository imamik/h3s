// Package ssh contains the functionality for executing SSH commands
package ssh

import (
	"fmt"
	"h3s/internal/utils/ip"
	"log"
	"strings"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"golang.org/x/crypto/ssh"
)

// ExecuteWithSsh executes a command on a remote server using SSH
func ExecuteWithSsh(
	privateSshKeyPath string,
	remote *hcloud.Server,
	command string,
) (string, error) {
	remoteIp := ip.FirstAvailable(remote)
	if err := removeKnownHostsEntry(remoteIp); err != nil {
		return "", err
	}

	// SSH client configuration
	sshConfig, err := getConfig(privateSshKeyPath)
	if err != nil {
		return "", fmt.Errorf("unable to create SSH config: %w", err)
	}

	// Connect to the remote server
	client, err := dialWithRetries(remoteIp, sshConfig, 5*time.Second, 5)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}

	// Run the command
	res, err := run(client, command)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(res), nil
}

// ExecuteViaProxy executes a command on a remote server using a gateway server
func ExecuteViaProxy(
	privateSshKeyPath string,
	gateway *hcloud.Server,
	remote *hcloud.Server,
	command string,
) (string, error) {
	// If the gateway is nil, execute the command directly on the remote server
	if gateway == nil {
		return ExecuteWithSsh(privateSshKeyPath, remote, command)
	}

	// Find the first available IP address of the gateway and the remote server
	proxyIp := ip.FirstAvailable(gateway)
	remoteIp := ip.FirstAvailable(remote)
	if err := removeKnownHostsEntry(proxyIp); err != nil {
		return "", err
	}

	// SSH client configuration
	sshConfig, err := getConfig(privateSshKeyPath)
	if err != nil {
		return "", fmt.Errorf("unable to create SSH config: %w", err)
	}

	// Connect to the gateway
	proxyConn, err := dialWithRetries(proxyIp, sshConfig, 5*time.Second, 5)
	if err != nil {
		return "", fmt.Errorf("unable to connect to gateway: %w", err)
	}
	defer func(proxyConn *ssh.Client) {
		if closeErr := proxyConn.Close(); closeErr != nil {
			fmt.Printf("unable to close gateway connection: %v\n", closeErr)
		}
	}(proxyConn)

	// Create a tunnel to the remote server
	tunnel, err := proxyConn.Dial("tcp", remoteIp+":22")
	if err != nil {
		return "", fmt.Errorf("unable to create tunnel to remote: %w", err)
	}

	// Connect to the remote server through the tunnel
	clientConn, chans, reqs, err := ssh.NewClientConn(tunnel, remoteIp+":22", sshConfig)
	if err != nil {
		return "", fmt.Errorf("unable to establish connection to remote: %w", err)
	}
	client := ssh.NewClient(clientConn, chans, reqs)

	// Run the command
	res, err := run(client, command)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(res), nil
}
