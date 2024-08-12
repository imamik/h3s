package k8s

import (
	"fmt"
	"h3s/internal/cluster"
	"h3s/internal/utils/common"
	"strings"
)

// Token creates a token for a user in a namespace
func Token(ctx *cluster.Cluster, namespace string, user string, hours int) (string, error) {
	// Create the command to execute on the remote server
	duration := fmt.Sprintf("%dh", hours)
	cmd := fmt.Sprintf("kubectl -n %s create token %s --duration=%s", namespace, user, duration)

	// Execute the command on the remote server
	bearer, err := common.SSH(ctx, cmd)
	if err != nil {
		return "", err
	}

	// Ensure the token is a single line string and return
	bearer = strings.ReplaceAll(bearer, "\n", "")
	return bearer, nil
}
