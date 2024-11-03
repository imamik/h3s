// Package execute contains the functionality for executing commands locally
package execute

import "os/exec"

// Local executes a command locally
func Local(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
