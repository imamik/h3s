// Package ssh contains utilities for executing commands over SSH
package ssh

import (
	"fmt"
	"h3s/internal/utils/ip"
	"strings"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"golang.org/x/crypto/ssh"
)

// ExecuteWithSSH executes a command on a remote server using SSH
func ExecuteWithSSH(privateSSHKeyPath string, remote *hcloud.Server, command string) (string, error) {
	if remote == nil {
		return "", fmt.Errorf("remote server is nil")
	}
	remoteIP := ip.FirstAvailable(remote)
	if err := removeKnownHostsEntry(remoteIP); err != nil {
		return "", err
	}

	// SSH client configuration
	sshConfig, err := getConfig(privateSSHKeyPath)
	if err != nil {
		return "", fmt.Errorf("unable to create SSH config: %w", err)
	}

	// Connect to the remote server
	client, err := dialWithRetries(remoteIP, sshConfig, 5*time.Second, 5)
	if err != nil {
		return "", fmt.Errorf("failed to dial: %w", err)
	}

	// Run the command
	res, err := run(client, command)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(res), nil
}

// ExecuteViaProxy executes a command on a remote server using a gateway server
func ExecuteViaProxy(privateSSHKeyPath string, gateway, remote *hcloud.Server, command string) (string, error) {
	// If the gateway is nil, execute the command directly on the remote server
	if gateway == nil {
		return ExecuteWithSSH(privateSSHKeyPath, remote, command)
	}

	// Find the first available IP address of the gateway and the remote server
	proxyIP := ip.FirstAvailable(gateway)
	remoteIP := ip.FirstAvailable(remote)
	if err := removeKnownHostsEntry(proxyIP); err != nil {
		return "", err
	}

	// SSH client configuration
	sshConfig, err := getConfig(privateSSHKeyPath)
	if err != nil {
		return "", fmt.Errorf("unable to create SSH config: %w", err)
	}

	// Connect to the gateway
	proxyConn, err := dialWithRetries(proxyIP, sshConfig, 5*time.Second, 5)
	if err != nil {
		return "", fmt.Errorf("unable to connect to gateway: %w", err)
	}
	defer func(proxyConn *ssh.Client) {
		if closeErr := proxyConn.Close(); closeErr != nil {
			fmt.Printf("unable to close gateway connection: %v\n", closeErr)
		}
	}(proxyConn)

	// Create a tunnel to the remote server
	tunnel, err := proxyConn.Dial("tcp", remoteIP+":22")
	if err != nil {
		return "", fmt.Errorf("unable to create tunnel to remote: %w", err)
	}

	// Connect to the remote server through the tunnel
	clientConn, chans, reqs, err := ssh.NewClientConn(tunnel, remoteIP+":22", sshConfig)
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
