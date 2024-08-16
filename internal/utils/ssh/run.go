package ssh

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"h3s/internal/utils/logger"
	"io"
	"time"
)

// dialWithRetries dials the SSH server with retries
func dialWithRetries(ip string, sshConfig *ssh.ClientConfig, retryInterval time.Duration, maxRetries int) (*ssh.Client, error) {
	for i := 0; i < maxRetries; i++ {
		c, err := ssh.Dial("tcp", ip+":22", sshConfig)
		if err == nil && c != nil {
			return c, err
		}
		retryBackoff := time.Duration(i+1) * retryInterval
		logger.LogResourceEvent(logger.Server, "SSH", ip, logger.Failure, fmt.Sprintf("Failed to dial: %s, retrying in %s", err, retryBackoff))
		time.Sleep(retryBackoff)
	}
	return nil, fmt.Errorf("failed to dial %s after %d retries", ip, maxRetries)
}

func run(client *ssh.Client, command string) (string, error) {

	defer func(client *ssh.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("unable to close client connection: %v\n", err)
		}
	}(client)

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

	var output = ""
	stdoutPipe, err := session.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("unable to create stdout pipe: %w", err)
	}

	var errorOutput = ""
	stderrPipe, err := session.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("unable to create stderr pipe: %w", err)
	}

	fmt.Printf("===============================================================================================\n")

	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			output += scanner.Text() + "\n"
			fmt.Println(scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			errorOutput += scanner.Text()
			fmt.Println(scanner.Text())
		}
	}()

	if err := session.Run(command); err != nil {
		return "", fmt.Errorf("command execution failed: %w, %s", err, errorOutput)
	}

	fmt.Printf("===============================================================================================\n")

	return output, nil
}
