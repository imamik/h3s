package kubectl

import (
	"fmt"
	"os/exec"
	"strings"
)

func Kubectl(args []string) {
	args = append(args, "--kubeconfig=k3s.yaml")

	// Joining all arguments to form a command string
	cmdStr := strings.Join(args, " ")

	// Executing the kubectl command with the provided arguments
	output, err := exec.Command("kubectl", args...).CombinedOutput()

	// Handling the command execution result
	if err != nil {
		fmt.Printf("Error executing kubectl with args [%s]: %s\n", cmdStr, err)
		fmt.Println(string(output)) // Printing any output that might have been generated
		return
	}

	// If successful, print the output
	fmt.Println(string(output))
}
