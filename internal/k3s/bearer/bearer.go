package bearer

import (
	"fmt"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/ssh"
)

func GetBearerToken(ctx clustercontext.ClusterContext, namespace string, user string, hours int) (string, error) {
	duration := fmt.Sprintf("%dh", hours)
	cmd := fmt.Sprintf("kubectl -n %s create token %s --duration=%s", namespace, user, duration)
	bearer, err := ssh.SSH(ctx, cmd)
	if err != nil {
		return "", err
	}
	return bearer, nil
}
