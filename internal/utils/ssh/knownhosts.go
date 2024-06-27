package ssh

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func removeKnownHostsEntry(ip string) {
	knownHostsFile := os.Getenv("HOME") + "/.ssh/known_hosts"
	tempFile := os.Getenv("HOME") + "/.ssh/temp_known_hosts"

	// Open original known_hosts file
	in, err := os.Open(knownHostsFile)
	if err != nil {
		log.Fatalf("Failed to open known_hosts file: %s", err)
	}
	defer in.Close()

	// Create a temporary file
	out, err := os.Create(tempFile)
	if err != nil {
		log.Fatalf("Failed to create temporary file: %s", err)
	}
	defer out.Close()

	// Create a scanner to read the known_hosts file line by line
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, ip) {
			// If the line does not contain the IP, write it to the temporary file
			_, err := out.WriteString(line + "\n")
			if err != nil {
				log.Fatalf("Failed to write to temporary file: %s", err)
			}
		}
	}

	// Replace the original known_hosts file with the temporary file
	err = os.Rename(tempFile, knownHostsFile)
	if err != nil {
		log.Fatalf("Failed to replace known_hosts file: %s", err)
	}
}
