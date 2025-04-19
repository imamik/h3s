package ssh

import (
	"bufio"
	"fmt"
	"h3s/internal/utils/logger"
	"io"
	"time"

	"golang.org/x/crypto/ssh"
)

// dialWithRetries dials the SSH server with retries
func dialWithRetries(ip string, sshConfig *ssh.ClientConfig, retryInterval time.Duration, maxRetries int) (*ssh.Client, error) {
	l := logger.New(nil, logger.Server, "Dial", ip)
	defer l.LogEvents()

	for i := 0; i < maxRetries; i++ {
		c, err := ssh.Dial("tcp", ip+":22", sshConfig)
		if err == nil && c != nil {
			l.AddEvent(logger.Success)
			return c, nil
		}
		retryBackoff := time.Duration(i+1) * retryInterval
		l.AddEvent(logger.Failure, err)
		l.AddEvent(logger.Failure, fmt.Errorf("retrying SSH dial in %s: %w", ip, err))
		time.Sleep(retryBackoff)
	}

	errMsg := fmt.Sprintf("Failed to dial %s after %d retries", ip, maxRetries)
	l.AddEvent(logger.Failure, fmt.Errorf(errMsg))
	return nil, fmt.Errorf("%s", errMsg)
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
		closeErr := session.Close()
		if closeErr != nil && closeErr != io.EOF {
			fmt.Printf("unable to close session: %v\n", closeErr)
		}
	}(session)

	output := ""
	stdoutPipe, err := session.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("unable to create stdout pipe: %w", err)
	}

	errorOutput := ""
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
