package ssh

import (
	"bufio"
	"fmt"
	"h3s/internal/utils/file"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

const user = "root"

// getConfig returns the SSH client configuration
func getConfig(privateSSHKeyPath string) (*ssh.ClientConfig, error) {
	// Load private key
	key, err := file.New(privateSSHKeyPath).Load().GetBytes()
	if err != nil {
		return nil, fmt.Errorf("unable to load private key: %w", err)
	}

	// Parse private key
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %w", err)
	}

	// SSH client configuration
	return &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		// #nosec G106
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // ignore host key verification
	}, nil
}

// removeKnownHostsEntry removes the known_hosts entry for the given IP
func removeKnownHostsEntry(ip string) error {
	homeDir, err := os.UserHomeDir() // Safely get the home directory
	if err != nil {
		return err
	}
	knownHostsFile := homeDir + "/.ssh/known_hosts"
	tempFile := homeDir + "/.ssh/temp_known_hosts"

	// Open original known_hosts file
	// #nosec G304
	in, err := os.Open(knownHostsFile)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := in.Close(); closeErr != nil {
			// Handle the error from closing the file
			err = closeErr
		}
	}()

	// Create a temporary file
	// #nosec G304
	out, err := os.Create(tempFile)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := out.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	// Create a scanner to read the known_hosts file line by line
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, ip) {
			// If the line does not contain the IP, write it to the temporary file
			_, err := out.WriteString(line + "\n")
			if err != nil {
				return err
			}
		}
	}

	// Replace the original known_hosts file with the temporary file
	return os.Rename(tempFile, knownHostsFile)
}
