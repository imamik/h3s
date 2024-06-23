package ssh

import (
	"bufio"
	"fmt"
	"github.com/melbahja/goph"
	"log"
)

func run(client *goph.Client, command string) {

	// Create a new command
	cmd, err := client.Command(command)
	if err != nil {
		log.Fatal(err)
	}

	// Get the command's stdout and stderr pipes
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	// Stream the stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Println("Error reading stdout:", err)
		}
	}()

	// Stream the stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Println("Error reading stderr:", err)
		}
	}()

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}
