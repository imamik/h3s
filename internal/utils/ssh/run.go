package ssh

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
)

func Run(client *ssh.Client, command string) (string, error) {

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
