package k8s

import (
	"fmt"
	"h3s/internal/cluster"
	"h3s/internal/utils/common"
	"h3s/internal/utils/kubectl"
	"strings"
)

// Token creates a token for a user in a namespace
func Token(ctx *cluster.Cluster, namespace, user string, hours int) (string, error) {
	// Create the command to execute on the remote server
	duration := fmt.Sprintf("%dh", hours)
	cmd, err := kubectl.
		New().
		Namespace(namespace).
		AddArgs("create", "token", user, fmt.Sprintf("--duration=%s", duration)).
		String()
	if err != nil {
		return "", err
	}

	// Execute the command on the remote server
	bearer, err := common.SSH(ctx, cmd)
	if err != nil {
		return "", err
	}

	// Ensure the token is a single line string and return
	bearer = strings.ReplaceAll(bearer, "\n", "")
	return bearer, nil
}
