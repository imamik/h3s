package ssh

import (
	"bufio"
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"golang.org/x/crypto/ssh"
	"hcloud-k3s-cli/internal/clustercontext"
	"io"
)

func Execute(
	ctx clustercontext.ClusterContext,
	proxy *hcloud.Server,
	remote *hcloud.Server,
	command string,
) (string, error) {
	user := "root"
	proxyIp := proxy.PublicNet.IPv4.IP.String()
	remoteIp := remote.PrivateNet[0].IP.String()

	// Load private key
	key, err := ReadPrivateKeyFromFile(ctx)

	signer, err := ssh.ParsePrivateKey([]byte(key))
	if err != nil {
		return "", fmt.Errorf("unable to parse private key: %w", err)
	}

	// SSH client configuration
	sshConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to proxy
	proxyConn, err := ssh.Dial("tcp", proxyIp+":22", sshConfig)
	if err != nil {
		return "", fmt.Errorf("unable to connect to proxy: %w", err)
	}
	defer func(proxyConn *ssh.Client) {
		err := proxyConn.Close()
		if err != nil {
			fmt.Printf("unable to close proxy connection: %v\n", err)
		}
	}(proxyConn)

	// Create a tunnel to the remote server via the proxy
	tunnel, err := proxyConn.Dial("tcp", remoteIp+":22")
	if err != nil {
		return "", fmt.Errorf("unable to create tunnel to remote: %w", err)
	}

	// Establish SSH connection to remote server through the tunnel
	clientConn, chans, reqs, err := ssh.NewClientConn(tunnel, remoteIp+":22", sshConfig)
	if err != nil {
		return "", fmt.Errorf("unable to establish connection to remote: %w", err)
	}
	client := ssh.NewClient(clientConn, chans, reqs)
	defer func(client *ssh.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("unable to close client connection: %v\n", err)
		}
	}(client)

	// Run command on remote server
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("unable to create session on remote: %w", err)
	}
	defer func(session *ssh.Session) {
		err := session.Close()
		if err != nil && err != io.EOF {
			fmt.Printf("unable to close session: %v\n", err)
		}
	}(session)

	stdoutPipe, err := session.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("unable to create stdout pipe: %w", err)
	}

	stderrPipe, err := session.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("unable to create stderr pipe: %w", err)
	}

	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	if err := session.Run(command); err != nil {
		return "", fmt.Errorf("command execution failed: %w", err)
	}

	return "", nil
}
